package k8sutil

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func CreateService(client clientv1.ServiceInterface, service *v1.Service) (*v1.Service, error) {
	created, err := client.Create(service)

	if err != nil {
		return nil, fmt.Errorf("creating service failed")
	}

	return created, nil
}
