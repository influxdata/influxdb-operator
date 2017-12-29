package k8sutil

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PortInput struct {
	Name     string
	Port     int32
	Protocol int32
}

type DeploymentInput struct {
	Namespace   string
	Name        string
	Image       string
	Replicas    *int32
	Labels      map[string]string
	Annotations map[string]string
	EnvVars     map[string]string
	Ports       []PortInput
}

// IncrementDeployment Create if it is the first one or Update a deployment.
func (k *K8sWrap) IncrementDeployment(input DeploymentInput) error {
	deployment, err := k.kubeClient.ExtensionsV1beta1().Deployments(input.Namespace).Get(input.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if len(deployment.Name) == 0 {
		k.kubeClient.ExtensionsV1beta1().Deployments(input.Namespace).Create(deployment)
	} else {
		if *deployment.Spec.Replicas != *input.Replicas {
			deployment.Spec.Replicas = input.Replicas
			if _, err := k.kubeClient.ExtensionsV1beta1().Deployments(input.Namespace).Update(deployment); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

// DeleteDeployments deletes deployments and replica sets filtering via labels
func (k *K8sWrap) DeleteDeploymentsSelectedByLabels(namespace string, labels map[string]string) error {
	labelSelector := ""
	for k, v := range labels {
		labelSelector = fmt.Sprintf("%s,%s=%s", labelSelector, k, v)
	}
	deployments, err := k.kubeClient.ExtensionsV1beta1().Deployments(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return err
	}
	for _, deployment := range deployments.Items {
		deployment.Spec.Replicas = new(int32)
		deployment, err := k.kubeClient.ExtensionsV1beta1().Deployments(namespace).Update(&deployment)
		if err != nil {
			log.Printf("Could not scale deployment: %s ", deployment.Name)
		} else {
			log.Printf("Scaled deployment: %s to zero", deployment.Name)
		}
		err = k.kubeClient.ExtensionsV1beta1().Deployments(namespace).Delete(deployment.Name, &metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Could not delete deployments: %s ", deployment.Name)
		} else {
			log.Printf("Deleted deployment: %s", deployment.Name)
		}
	}

	replicaSets, err := k.kubeClient.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return err
	}
	for _, replicaSet := range replicaSets.Items {
		err := k.kubeClient.ExtensionsV1beta1().ReplicaSets(namespace).Delete(replicaSet.Name, &metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Could not delete replica sets: %s ", replicaSet.Name)
		} else {
			log.Printf("Deleted replica set: %s", replicaSet.Name)
		}
	}

	return nil
}
