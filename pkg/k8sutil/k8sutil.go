package k8sutil

import (
	"fmt"

	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func CreateService(client clientv1.ServiceInterface, service *v1.Service) error {
	_, err := client.Create(service)

	if apierrors.IsAlreadyExists(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("creating service failed: %s", err.Error())
	}

	return nil
}

func DeleteServices(client clientv1.ServiceInterface, serviceName string) error {
	return client.Delete(serviceName, &meta_v1.DeleteOptions{})
}

func CreateConfigMap(client clientv1.ConfigMapInterface, cm *v1.ConfigMap) error {
	_, err := client.Create(cm)
	if apierrors.IsAlreadyExists(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("creating config map failed: %s", err.Error())
	}

	return nil
}

func DeleteConfigMap(client clientv1.ConfigMapInterface, cm string) error {
	return client.Delete(cm, &meta_v1.DeleteOptions{})
}

func CreateDeployment(client clientv1beta1.DeploymentInterface, deployment *v1beta1.Deployment) error {
	_, err := client.Create(deployment)
	if apierrors.IsAlreadyExists(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

type DeploymentInput struct {
	Name            string
	Image           string
	ImagePullPolicy v1.PullPolicy
	Labels          map[string]string
	Selector        *meta_v1.LabelSelector
	Replicas        *int32
	Namespace       string
	Ports           []v1.ContainerPort
	VolumeMounts    []v1.VolumeMount
	Resources       v1.ResourceRequirements
	Volumes         []v1.Volume
	Envs            []v1.EnvVar
}

func NewDeployment(i DeploymentInput) *v1beta1.Deployment {
	return &v1beta1.Deployment{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      i.Name,
			Labels:    i.Labels,
			Namespace: i.Namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Selector: i.Selector,
			Replicas: i.Replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels:    i.Labels,
					Namespace: i.Namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:            i.Name,
							Image:           i.Image,
							ImagePullPolicy: i.ImagePullPolicy,
							Env:             i.Envs,
							Ports:           i.Ports,
							VolumeMounts:    i.VolumeMounts,
							Resources:       i.Resources,
						},
					},
					Volumes: i.Volumes,
				},
			},
		},
	}
}
