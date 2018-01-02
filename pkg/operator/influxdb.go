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
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	deploymentName := fmt.Sprintf("%s-%s", v1alpha1.InfluxDBKind, oret.GetName())
	policy := metav1.DeletePropagationForeground
	err := o.kubeClient.ExtensionsV1beta1().Deployments(oret.GetNamespace()).Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		log.Printf("Error deleting deployment %s. %s", deploymentName, err)
	}

	err = k8sutil.DeleteServices(o.kubeClient.CoreV1().Services(oret.GetNamespace()), deploymentName)

	if err != nil {
		log.Printf("Error deleting deployment service: %s. %s", deploymentName, err)
	}

	err = k8sutil.DeleteConfigMap(o.kubeClient.CoreV1().ConfigMaps(oret.GetNamespace()), deploymentName)
	if err != nil {
		log.Printf("Error deleting config map: %s. %s", deploymentName, err)
	}
}

func (o *Operator) handleAddInfluxDB(obj interface{}) {
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	influxdbSpec := obj.(*v1alpha1.Influxdb)

	labels := o.config.Labels
	for k, v := range influxdbSpec.GetLabels() {
		labels[k] = v
	}

	choosenName := fmt.Sprintf("%s-%s", v1alpha1.InfluxDBKind, oret.GetName())
	labels["name"] = choosenName
	labels["resource"] = v1alpha1.InfluxDBKind

	config := createInfluxDBConfiguration(influxdbSpec.Spec)
	cm, err := makeInfluxDBConfigMap(choosenName, config)

	if err != nil {
		log.Print(err)
		return
	}

	err = k8sutil.CreateConfigMap(o.kubeClient.Core().ConfigMaps(oret.GetNamespace()), cm)

	if err != nil {
		log.Print(err)
		return
	}

	deployment := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      labels["name"],
			Labels:    labels,
			Namespace: oret.GetNamespace(),
		},
		Spec: v1beta1.DeploymentSpec{
			Selector: influxdbSpec.Spec.Selector,
			Replicas: influxdbSpec.Spec.Replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    labels,
					Namespace: oret.GetNamespace(),
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:            oret.GetName(),
							Image:           influxdbSpec.Spec.BaseImage,
							ImagePullPolicy: "Always",
							Env:             []v1.EnvVar{},
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
							//Resources:    v1.ResourceRequirements{},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "influxdb-config",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: choosenName,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = o.kubeClient.AppsV1beta1().Deployments(oret.GetNamespace()).Create(deployment)
	if err != nil {
		log.Print(err)
	}

	svc := makeInfluxDBService(choosenName, o.config)
	err = k8sutil.CreateService(o.kubeClient.CoreV1().Services(oret.GetNamespace()), svc)

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
