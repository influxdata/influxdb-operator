package v1alpha1

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	InfluxDBKind   = "influxdb"
	InfluxDBPlural = "influxdbs"
)

type influxdbs struct {
	restClient rest.Interface
	client     dynamic.ResourceInterface
	crdKind    CrdKind
	namespace  string
}

func (i *influxdbs) List(opts metav1.ListOptions) (runtime.Object, error) {
	panic("not implemented")
}

func (i *influxdbs) Get(name string, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *influxdbs) Delete(name string, opts *metav1.DeleteOptions) error {
	panic("not implemented")
}

func (i *influxdbs) DeleteCollection(deleteOptions *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	panic("not implemented")
}

func (i *influxdbs) Create(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *influxdbs) Update(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *influxdbs) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	panic("not implemented")
}

func (i *influxdbs) Patch(name string, pt types.PatchType, data []byte) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func newInfluxdbs(r rest.Interface, c *dynamic.Client, namespace string) *influxdbs {
	return &influxdbs{
		restClient: r,
		client: c.Resource(
			&metav1.APIResource{
				Kind:       InfluxDBKind,
				Name:       InfluxDBPlural,
				Namespaced: true,
			},
			namespace,
		),
		namespace: namespace,
	}
}

func UnstructuredFromInfluxDB(p *Influxdb) (*unstructured.Unstructured, error) {
	p.TypeMeta.Kind = InfluxDBKind
	// TODO: Naaah It's not right.
	p.TypeMeta.APIVersion = "gianarb.com/v1alpha1"
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	var r unstructured.Unstructured
	if err := json.Unmarshal(b, &r.Object); err != nil {
		return nil, err
	}
	return &r, nil
}
