package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	Group   = "tick.gianarb.com"
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

var DefaultCrdKinds CrdKinds = CrdKinds{
	KindsString: "",
	InfluxDB:    CrdKind{Plural: InfluxDBPlural, Kind: InfluxDBKind},
}

type TickV1alpha1Interface interface {
	RESTClient() rest.Interface
}

type TickV1alpha1Client struct {
	restClient    rest.Interface
	dynamicClient *dynamic.Client
}

func (c *TickV1alpha1Client) InfluxDBs(namespace string) dynamic.ResourceInterface {
	return newInfluxdbs(c.restClient, c.dynamicClient, namespace)
}

func NewForConfig(c *rest.Config) (*rest.RESTClient, dynamic.Interface, error) {
	config := *c
	config.GroupVersion = &schema.GroupVersion{
		Group:   Group,
		Version: Version,
	}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewClient(&config)
	if err != nil {
		panic(err)
	}
	return client, dynamicClient, nil
}
