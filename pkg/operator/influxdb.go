package operator

import (
	"bytes"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
	"github.com/gianarb/influxdb-operator/pkg/k8sutil"
	"github.com/influxdata/influxdb/cmd/influxd/run"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
)

func registerInfluxInformer(operator *Operator) {
	operator.influxInformer = cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc:  operator.tickCs.InfluxDBs(metav1.NamespaceAll).List,
			WatchFunc: operator.tickCs.InfluxDBs(metav1.NamespaceAll).Watch,
		},
		&v1alpha1.Influxdb{}, 0, cache.Indexers{},
	)

	operator.influxInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    operator.handleAddInfluxDB,
		UpdateFunc: nil,
		DeleteFunc: operator.handleDeleteInfluxDB,
	})
}

func (o *Operator) handleDeleteInfluxDB(obj interface{}) {
	spec := obj.(*v1alpha1.Influxdb)
	deploymentName := fmt.Sprintf("%s-%s", v1alpha1.InfluxDBKind, spec.GetName())
	policy := metav1.DeletePropagationForeground
	err := o.kubeClient.ExtensionsV1beta1().Deployments(spec.GetNamespace()).Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		log.Printf("Error deleting deployment %s. %s", deploymentName, err)
	}

	err = k8sutil.DeleteServices(o.kubeClient.CoreV1().Services(spec.GetNamespace()), deploymentName)

	if err != nil {
		log.Printf("Error deleting deployment service: %s. %s", deploymentName, err)
	}

	err = k8sutil.DeleteConfigMap(o.kubeClient.CoreV1().ConfigMaps(spec.GetNamespace()), deploymentName)
	if err != nil {
		log.Printf("Error deleting config map: %s. %s", deploymentName, err)
	}
}

func makeInfluxDBDeployment(deploymentName string, influxdbSpec *v1alpha1.Influxdb) *v1beta1.Deployment {
	labels := map[string]string{}
	for k, v := range influxdbSpec.GetLabels() {
		labels[k] = v
	}

	labels["name"] = deploymentName
	labels["resource"] = v1alpha1.InfluxDBKind
	i := k8sutil.DeploymentInput{
		Name:            labels["name"],
		Image:           influxdbSpec.Spec.BaseImage,
		ImagePullPolicy: influxdbSpec.Spec.ImagePullPolicy,
		Labels:          labels,
		Selector:        influxdbSpec.Spec.Selector,
		Replicas:        influxdbSpec.Spec.Replicas,
		Namespace:       influxdbSpec.GetNamespace(),
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				Name:          "http",
				ContainerPort: 8086,
				Protocol:      v1.ProtocolTCP,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "influxdb-config",
				SubPath:   "config.toml",
				MountPath: "/etc/influxdb/influxdb.conf",
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "influxdb-config",
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: deploymentName,
						},
					},
				},
			},
		},
	}
	return k8sutil.NewDeployment(i)
}

func (o *Operator) handleAddInfluxDB(obj interface{}) {
	influxdbSpec := obj.(*v1alpha1.Influxdb)

	choosenName := fmt.Sprintf("%s-%s", v1alpha1.InfluxDBKind, influxdbSpec.GetName())
	config := createInfluxDBConfiguration(influxdbSpec.Spec)
	cm, err := makeInfluxDBConfigMap(choosenName, config)

	if err != nil {
		log.Print(err)
		return
	}

	err = k8sutil.CreateConfigMap(o.kubeClient.Core().ConfigMaps(influxdbSpec.GetNamespace()), cm)

	if err != nil {
		log.Print(err)
		return
	}

	deployment := makeInfluxDBDeployment(choosenName, influxdbSpec)
	err = k8sutil.CreateDeployment(o.kubeClient.AppsV1beta1().Deployments(influxdbSpec.GetNamespace()), deployment)
	if err != nil {
		log.Print(err)
	}

	svc := makeInfluxDBService(choosenName, o.config)
	err = k8sutil.CreateService(o.kubeClient.CoreV1().Services(influxdbSpec.GetNamespace()), svc)

	if err != nil {
		log.Print(err)
		return
	}
}

func makeInfluxDBService(name string, config Config) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Labels:          config.Labels,
			OwnerReferences: nil,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "influx",
					Port:       8086,
					TargetPort: intstr.FromInt(8086),
					Protocol:   v1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"name": name,
			},
		},
	}
}

func createInfluxDBConfiguration(spec v1alpha1.InfluxdbSpec) *run.Config {
	config := run.NewConfig()
	config.Meta.Dir = "/var/lib/influxdb/meta"
	config.Data.Dir = "/var/lib/influxdb/data"
	config.Data.WALDir = "/var/lib/influxdb/wal"
	return config
}

func makeInfluxDBConfigMap(name string, config *run.Config) (*v1.ConfigMap, error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(&config); err != nil {
		return nil, err
	}

	data := map[string]string{
		"config.toml": buf.String(),
	}

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}, nil
}
