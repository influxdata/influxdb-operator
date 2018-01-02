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
		return fmt.Errorf("creating service failed: %s", err)
	}

	return nil
}

func DeleteServices(client clientv1.ServiceInterface, serviceName string) error {
	return client.Delete(serviceName, &meta_v1.DeleteOptions{})
}
