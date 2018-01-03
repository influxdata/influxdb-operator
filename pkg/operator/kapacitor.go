package operator

import (
	"bytes"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
	"github.com/gianarb/influxdb-operator/pkg/k8sutil"
	"github.com/influxdata/kapacitor/server"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
)

func registerKapacitorInformer(operator *Operator) {
	operator.kapacitorInformer = cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc:  operator.tickCs.Kapacitors(metav1.NamespaceAll).List,
			WatchFunc: operator.tickCs.Kapacitors(metav1.NamespaceAll).Watch,
		},
		&v1alpha1.Kapacitor{}, 0, cache.Indexers{},
	)

	operator.kapacitorInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    operator.handleAddKapacitor,
		UpdateFunc: nil,
		DeleteFunc: operator.handleDeleteKapacitor,
	})
}

func createKapacitorConfiguration(spec *v1alpha1.Kapacitor) *server.Config {
	config, _ := server.NewDemoConfig()
	return config
}

func (o *Operator) handleDeleteKapacitor(obj interface{}) {
	spec := obj.(*v1alpha1.Kapacitor)

	deploymentName := fmt.Sprintf("%s-%s", v1alpha1.KapacitorKind, spec.GetName())
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

func makeKapacitorDeployment(deploymentName string, spec *v1alpha1.Kapacitor) *v1beta1.Deployment {
	labels := map[string]string{}
	for k, v := range spec.GetLabels() {
		labels[k] = v
	}
	labels["name"] = deploymentName
	labels["resource"] = v1alpha1.KapacitorKind
	i := k8sutil.DeploymentInput{
		Name:            labels["name"],
		Image:           spec.Spec.BaseImage,
		ImagePullPolicy: spec.Spec.ImagePullPolicy,
		Labels:          labels,
		Selector:        spec.Spec.Selector,
		Replicas:        spec.Spec.Replicas,
		Namespace:       spec.GetNamespace(),
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				Name:          "http",
				ContainerPort: 9092,
				Protocol:      v1.ProtocolTCP,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "kapacitor-config",
				SubPath:   "config.toml",
				MountPath: "/etc/influxdb/influxdb.conf",
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "kapacitor-config",
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

func makeKapacitorConfigMap(name string, config *server.Config) (*v1.ConfigMap, error) {
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

func (o *Operator) handleAddKapacitor(obj interface{}) {
	spec := obj.(*v1alpha1.Kapacitor)

	choosenName := fmt.Sprintf("%s-%s", v1alpha1.KapacitorKind, spec.GetName())

	kapacitorConfig := createKapacitorConfiguration(spec)
	cm, err := makeKapacitorConfigMap(choosenName, kapacitorConfig)
	if err != nil {
		log.Print(err)
		return
	}
	err = k8sutil.CreateConfigMap(o.kubeClient.Core().ConfigMaps(spec.GetNamespace()), cm)
	if err != nil {
		log.Print(err)
		return
	}

	deployment := makeKapacitorDeployment(choosenName, spec)
	err = k8sutil.CreateDeployment(o.kubeClient.AppsV1beta1().Deployments(spec.GetNamespace()), deployment)
	if err != nil {
		log.Print(err)
		return
	}

	svc := makeKapacitorService(choosenName, o.config)
	err = k8sutil.CreateService(o.kubeClient.CoreV1().Services(spec.GetNamespace()), svc)

	if err != nil {
		log.Print(err)
		return
	}
}

func makeKapacitorService(name string, config Config) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Labels:          config.Labels,
			OwnerReferences: nil,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "kapacitor",
					Port:       9092,
					TargetPort: intstr.FromInt(9092),
					Protocol:   v1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"name": name,
			},
		},
	}
}
