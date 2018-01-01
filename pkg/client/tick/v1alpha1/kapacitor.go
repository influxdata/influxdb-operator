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
	KapacitorKind       = "kapacitor"
	KapacitorPlural     = "kapacitors"
	KapacitorAPIVersion = "gianarb.com/v1alpha1" //TODO(gianarb): decide what to do with this
)

type kapacitors struct {
	restClient rest.Interface
	client     dynamic.ResourceInterface
	crdKind    CrdKind
	namespace  string
}

func (i *kapacitors) List(opts metav1.ListOptions) (runtime.Object, error) {
	req := i.restClient.Get().Namespace(i.namespace).Resource(KapacitorPlural)

	buf, err := req.DoRaw()

	if err != nil {
		return nil, err
	}

	var list KapacitorList
	return &list, json.Unmarshal(buf, &list)
}

func (i *kapacitors) Get(name string, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	cur, err := i.client.Get(name, opts)

	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (i *kapacitors) Delete(name string, opts *metav1.DeleteOptions) error {
	log.Printf("Delete: %s %v", name, opts)
	return nil
}

func (i *kapacitors) DeleteCollection(deleteOptions *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	log.Printf("DeleteCollection")
	return nil
}

func (i *kapacitors) Create(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *kapacitors) Update(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *kapacitors) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	r, err := i.restClient.Get().
		Prefix("watch").
		Namespace(i.namespace).
		Resource("kapacitors").
		//TODO: crdKind is not populated
		//Resource(i.crdKind.Plural).
		Stream()
	if err != nil {
		return nil, err
	}
	return watch.NewStreamWatcher(&kapacitorDecoder{
		dec:   json.NewDecoder(r),
		close: r.Close,
	}), nil

}

func (i *kapacitors) Patch(name string, pt types.PatchType, data []byte) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func newKapacitors(r rest.Interface, c *dynamic.Client, namespace string) *kapacitors {
	return &kapacitors{
		restClient: r,
		client: c.Resource(
			&metav1.APIResource{
				Kind:       KapacitorKind,
				Name:       KapacitorPlural,
				Namespaced: true,
			},
			namespace,
		),
		namespace: namespace,
	}
}

type kapacitorDecoder struct {
	dec   *json.Decoder
	close func() error
}

func (j *kapacitorDecoder) Decode() (action watch.EventType, object runtime.Object, err error) {
	var e struct {
		Type   watch.EventType
		Object Kapacitor
	}
	if err := j.dec.Decode(&e); err != nil {
		return watch.Error, nil, err
	}
	return e.Type, &e.Object, nil
}

func (j *kapacitorDecoder) Close() {
	j.close()
}
