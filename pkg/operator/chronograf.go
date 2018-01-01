package operator

import (
	"fmt"
	"log"

	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	deploymentName := fmt.Sprintf("%s-%s", v1alpha1.ChronografPlural, oret.GetName())
	policy := metav1.DeletePropagationForeground
	err := o.kubeClient.ExtensionsV1beta1().Deployments(oret.GetNamespace()).Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		log.Printf("Error deleting deployment %s. %s", deploymentName, err)
	}
}

func (o *Operator) handleAddChronograf(obj interface{}) {
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	chronografSpec := obj.(*v1alpha1.Chronograf)

	labels := o.config.Labels
	labels["name"] = fmt.Sprintf("%s-%s", v1alpha1.ChronografPlural, oret.GetName())
	labels["resource"] = v1alpha1.ChronografPlural

	replicas := int32(1)
	deployment := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      labels["name"],
			Labels:    labels,
			Namespace: oret.GetNamespace(),
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    labels,
					Namespace: oret.GetNamespace(),
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:            oret.GetName(),
							Image:           chronografSpec.Spec.BaseImage,
							ImagePullPolicy: "Always",
							Env:             []v1.EnvVar{},
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									Name:          "http",
									ContainerPort: 8888,
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
	fmt.Print(err) //TODO: handle in a different way?
}
