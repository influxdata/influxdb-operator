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
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	deploymentName := fmt.Sprintf("%s-%s", v1alpha1.KapacitorKind, oret.GetName())
	policy := metav1.DeletePropagationForeground
	err := o.kubeClient.ExtensionsV1beta1().Deployments(oret.GetNamespace()).Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		log.Printf("Error deleting deployment %s. %s", deploymentName, err)
	}
}

func (o *Operator) handleAddKapacitor(obj interface{}) {
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	kapacitorSpec := obj.(*v1alpha1.Kapacitor)

	labels := o.config.Labels
	for k, v := range kapacitorSpec.GetLabels() {
		labels[k] = v
	}
	choosenName := fmt.Sprintf("%s-%s", v1alpha1.KapacitorKind, oret.GetName())
	labels["name"] = choosenName
	labels["resource"] = v1alpha1.KapacitorKind

	deployment := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      labels["name"],
			Labels:    labels,
			Namespace: oret.GetNamespace(),
		},
		Spec: v1beta1.DeploymentSpec{
			Selector: kapacitorSpec.Spec.Selector,
			Replicas: kapacitorSpec.Spec.Replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    labels,
					Namespace: oret.GetNamespace(),
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:            oret.GetName(),
							Image:           kapacitorSpec.Spec.BaseImage,
							ImagePullPolicy: "Always",
							Env:             []v1.EnvVar{},
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									Name:          "http",
									ContainerPort: 9092,
									Protocol:      v1.ProtocolTCP,
								},
							},
							//VolumeMounts: []v1.VolumeMount{},
							//Resources:    v1.ResourceRequirements{},
						},
					},
					//Volumes: []v1.Volume{},
				},
			},
		},
	}

	_, err := o.kubeClient.AppsV1beta1().Deployments(oret.GetNamespace()).Create(deployment)
	if err != nil {
		log.Print(err)
	}

	svc := makeKapacitorService(choosenName, o.config)
	_, err = k8sutil.CreateService(o.kubeClient.CoreV1().Services(oret.GetNamespace()), svc)

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
