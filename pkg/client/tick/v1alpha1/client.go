package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	Group   = "gianarb.com"
	Version = "v1alpha1"
)

type CrdKind struct {
	Kind   string
	Plural string
}

type CrdKinds struct {
	KindsString string
	InfluxDB    CrdKind
}

type TickV1alpha1Interface interface {
	RESTClient() rest.Interface
	InfluxDBs(namespace string) dynamic.ResourceInterface
}

type TickV1alpha1Client struct {
	restClient    *rest.RESTClient
	dynamicClient *dynamic.Client
}

func (c *TickV1alpha1Client) InfluxDBs(namespace string) dynamic.ResourceInterface {
	return newInfluxdbs(c.restClient, c.dynamicClient, namespace)
}

func (c *TickV1alpha1Client) Kapacitors(namespace string) dynamic.ResourceInterface {
	return newKapacitors(c.restClient, c.dynamicClient, namespace)
}

func (c *TickV1alpha1Client) Chronografs(namespace string) dynamic.ResourceInterface {
	return newChronografs(c.restClient, c.dynamicClient, namespace)
}

func NewForConfig(c *rest.Config) (*TickV1alpha1Client, error) {
	config := *c
	config.GroupVersion = &schema.GroupVersion{
		Group:   Group,
		Version: Version,
	}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	dynamicClient, err := dynamic.NewClient(&config)

	if err != nil {
		return nil, err
	}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &TickV1alpha1Client{
		restClient:    client,
		dynamicClient: dynamicClient,
	}, nil
}
