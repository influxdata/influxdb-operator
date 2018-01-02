package k8sutil

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func CreateService(client clientv1.ServiceInterface, service *v1.Service) (*v1.Service, error) {
	//TODO(fntlnz): check if the service already exists and ignore it in case
	created, err := client.Create(service)

	if err != nil {
		return nil, fmt.Errorf("creating service failed")
	}

	return created, nil
}

func DeleteServices(client clientv1.ServiceInterface, serviceName string) error {
	return client.Delete(serviceName, &meta_v1.DeleteOptions{})
}
