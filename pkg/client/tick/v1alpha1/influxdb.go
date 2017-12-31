package v1alpha1

import (
	"encoding/json"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	InfluxDBKind       = "influxdb"
	InfluxDBPlural     = "influxdbs"
	InfluxDBApiVersion = "gianarb.com/v1alpha1" //TODO(gianarb): decide what to do with this
)

type influxdbs struct {
	restClient rest.Interface
	client     dynamic.ResourceInterface
	crdKind    CrdKind
	namespace  string
}

func (i *influxdbs) List(opts metav1.ListOptions) (runtime.Object, error) {
	log.Printf("Namespace: %s", i.namespace)   //TODO: fill this
	log.Printf("Plural: %s", i.crdKind.Plural) //TODO: fill this
	req := i.restClient.Get().Namespace("default").Resource(InfluxDBPlural)

	buf, err := req.DoRaw()

	if err != nil {
		return nil, err
	}

	var list InfluxdbList
	return &list, json.Unmarshal(buf, &list)
}

func (i *influxdbs) Get(name string, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	log.Printf("Get: %s %v", name, opts)

	cur, err := i.client.Get(name, opts)

	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (i *influxdbs) Delete(name string, opts *metav1.DeleteOptions) error {
	log.Printf("Delete: %s %v", name, opts)
	return nil
}

func (i *influxdbs) DeleteCollection(deleteOptions *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	log.Printf("DeleteCollection")
	return nil
}

func (i *influxdbs) Create(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *influxdbs) Update(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *influxdbs) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	r, err := i.restClient.Get().
		Prefix("watch").
		Namespace("default").
		Resource(i.crdKind.Plural).
		Stream()
	if err != nil {
		return nil, err
	}
	return watch.NewStreamWatcher(&influxdbDecoder{
		dec:   json.NewDecoder(r),
		close: r.Close,
	}), nil

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
	p.TypeMeta.APIVersion = InfluxDBApiVersion
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

func InfluxDBFromUnstructured(r *unstructured.Unstructured) (*Influxdb, error) {
	b, err := json.Marshal(r.Object)
	if err != nil {
		return nil, err
	}

	var i Influxdb
	if err := json.Unmarshal(b, &i); err != nil {
		return nil, err
	}

	i.TypeMeta.Kind = InfluxDBKind
	i.TypeMeta.APIVersion = InfluxDBApiVersion

	return &i, nil

}

type influxdbDecoder struct {
	dec   *json.Decoder
	close func() error
}

func (j *influxdbDecoder) Decode() (action watch.EventType, object runtime.Object, err error) {
	panic("not implemented")
}

func (j *influxdbDecoder) Close() {
	panic("not implemented")
}
