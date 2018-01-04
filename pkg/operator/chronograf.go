package operator

import (
	"fmt"
	"log"

	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
	"github.com/gianarb/influxdb-operator/pkg/k8sutil"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
)

func registerChronografInformer(operator *Operator) {
	operator.chronografInformer = cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc:  operator.tickCs.Chronografs(metav1.NamespaceAll).List,
			WatchFunc: operator.tickCs.Chronografs(metav1.NamespaceAll).Watch,
		},
		&v1alpha1.Chronograf{}, 0, cache.Indexers{},
	)

	operator.chronografInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    operator.handleAddChronograf,
		UpdateFunc: nil,
		DeleteFunc: operator.handleDeleteChronograf,
	})
}

func (o *Operator) handleDeleteChronograf(obj interface{}) {
	spec := obj.(*v1alpha1.Chronograf)
	deploymentName := fmt.Sprintf("%s-%s", v1alpha1.ChronografKind, spec.GetName())
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

}

func makeChronografDeployment(deploymentName string, spec *v1alpha1.Chronograf) *v1beta1.Deployment {
	labels := map[string]string{}
	for k, v := range spec.GetLabels() {
		labels[k] = v
	}
	labels["name"] = deploymentName
	labels["resource"] = v1alpha1.ChronografKind
	envs := []v1.EnvVar{
		{
			Value: fmt.Sprint(spec.Spec.Port),
			Name:  "PORT",
		},
		{
			Value: spec.Spec.Host,
			Name:  "HOST",
		},
		{
			Value: spec.Spec.BoltPath,
			Name:  "BOLT_PATH",
		},
		{
			Value: spec.Spec.CannedPath,
			Name:  "CANNED_PATH",
		},
		{
			Value: spec.Spec.AuthDuration,
			Name:  "AUTH_DURATION",
		},
		{
			Value: spec.Spec.LogLevel,
			Name:  "LOG_LEVEL",
		},
	}

	i := k8sutil.DeploymentInput{
		Name:            labels["name"],
		Image:           spec.Spec.BaseImage,
		ImagePullPolicy: spec.Spec.ImagePullPolicy,
		Labels:          labels,
		Selector:        spec.Spec.Selector,
		Replicas:        spec.Spec.Replicas,
		Envs:            envs,
		Namespace:       spec.GetNamespace(),
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				Name:          "http",
				ContainerPort: spec.Spec.Port,
				Protocol:      v1.ProtocolTCP,
			},
		},
	}
	return k8sutil.NewDeployment(i)
}

func setDefaultSpecValues(spec *v1alpha1.Chronograf) {
	if spec.Spec.Host == "" {
		spec.Spec.Host = "0.0.0.0"
	}
	if spec.Spec.Port == 0 {
		spec.Spec.Port = 8888
	}
	if spec.Spec.BoltPath == "" {
		spec.Spec.BoltPath = "/var/lib/chronograf/chronograf-v1-.db"
	}
	if spec.Spec.CannedPath == "" {
		spec.Spec.CannedPath = "/usr/share/chronograf/canned"
	}
	if spec.Spec.AuthDuration == "" {
		spec.Spec.AuthDuration = "720h"
	}

	if spec.Spec.LogLevel == "" {
		spec.Spec.AuthDuration = "info"
	}
}

func (o *Operator) handleAddChronograf(obj interface{}) {
	chronografSpec := obj.(*v1alpha1.Chronograf)
	setDefaultSpecValues(chronografSpec)
	choosenName := fmt.Sprintf("%s-%s", v1alpha1.ChronografKind, chronografSpec.GetName())

	deployment := makeChronografDeployment(choosenName, chronografSpec)
	err := k8sutil.CreateDeployment(o.kubeClient.AppsV1beta1().Deployments(chronografSpec.GetNamespace()), deployment)
	if err != nil {
		log.Print(err)
	}

	svc := makeChronografService(choosenName, o.config, chronografSpec)
	if err != nil {
		log.Print(err)
	}
	err = k8sutil.CreateService(o.kubeClient.CoreV1().Services(chronografSpec.GetNamespace()), svc)

	if err != nil {
		log.Print(err)
	}

}

func makeChronografService(name string, config Config, spec *v1alpha1.Chronograf) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Labels:          config.Labels,
			OwnerReferences: nil,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "ui",
					Port:       spec.Spec.Port,
					TargetPort: intstr.FromInt(int(spec.Spec.Port)),
					Protocol:   v1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"name": name,
			},
		},
	}
}
