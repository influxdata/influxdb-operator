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
	}
	return k8sutil.NewDeployment(i)
}

func (o *Operator) handleAddKapacitor(obj interface{}) {
	spec := obj.(*v1alpha1.Kapacitor)

	choosenName := fmt.Sprintf("%s-%s", v1alpha1.KapacitorKind, spec.GetName())
	deployment := makeKapacitorDeployment(choosenName, spec)
	err := k8sutil.CreateDeployment(o.kubeClient.AppsV1beta1().Deployments(spec.GetNamespace()), deployment)
	if err != nil {
		log.Print(err)
	}

	svc := makeKapacitorService(choosenName, o.config)
	err = k8sutil.CreateService(o.kubeClient.CoreV1().Services(spec.GetNamespace()), svc)

	if err != nil {
		log.Print(err)
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
