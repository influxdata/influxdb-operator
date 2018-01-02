package k8sutil

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
