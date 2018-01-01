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
	ChronografKind       = "chronograf"
	ChronografPlural     = "chronografs"
	ChronografAPIVersion = "gianarb.com/v1alpha1" //TODO(gianarb): decide what to do with this
)

type chronografs struct {
	restClient rest.Interface
	client     dynamic.ResourceInterface
	crdKind    CrdKind
	namespace  string
}

func (i *chronografs) List(opts metav1.ListOptions) (runtime.Object, error) {
	req := i.restClient.Get().Namespace(i.namespace).Resource(ChronografPlural)

	buf, err := req.DoRaw()

	if err != nil {
		return nil, err
	}

	var list ChronografList
	return &list, json.Unmarshal(buf, &list)
}

func (i *chronografs) Get(name string, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	cur, err := i.client.Get(name, opts)

	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (i *chronografs) Delete(name string, opts *metav1.DeleteOptions) error {
	log.Printf("Delete: %s %v", name, opts)
	return nil
}

func (i *chronografs) DeleteCollection(deleteOptions *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	log.Printf("DeleteCollection")
	return nil
}

func (i *chronografs) Create(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *chronografs) Update(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func (i *chronografs) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	r, err := i.restClient.Get().
		Prefix("watch").
		Namespace(i.namespace).
		Resource("chronografs").
		//TODO: crdKind is not populated
		//Resource(i.crdKind.Plural).
		Stream()
	if err != nil {
		return nil, err
	}
	return watch.NewStreamWatcher(&chronografDecoder{
		dec:   json.NewDecoder(r),
		close: r.Close,
	}), nil

}

func (i *chronografs) Patch(name string, pt types.PatchType, data []byte) (*unstructured.Unstructured, error) {
	panic("not implemented")
}

func newChronografs(r rest.Interface, c *dynamic.Client, namespace string) *chronografs {
	return &chronografs{
		restClient: r,
		client: c.Resource(
			&metav1.APIResource{
				Kind:       ChronografKind,
				Name:       ChronografPlural,
				Namespaced: true,
			},
			namespace,
		),
		namespace: namespace,
	}
}

type chronografDecoder struct {
	dec   *json.Decoder
	close func() error
}

func (j *chronografDecoder) Decode() (action watch.EventType, object runtime.Object, err error) {
	var e struct {
		Type   watch.EventType
		Object Chronograf
	}
	if err := j.dec.Decode(&e); err != nil {
		return watch.Error, nil, err
	}
	return e.Type, &e.Object, nil
}

func (j *chronografDecoder) Close() {
	j.close()
}
