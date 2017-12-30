package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Influxdb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              InfluxdbSpec `json:"spec"`
}

func (i *Influxdb) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *Influxdb) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type InfluxdbList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*Influxdb `json:"items"`
}

func (i *InfluxdbList) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *InfluxdbList) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type InfluxdbSpec struct {
	BaseImage string `json:"baseImage,omitempty"`
}

func (i *InfluxdbSpec) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *InfluxdbSpec) DeepCopyObject() runtime.Object {
	panic("not implemented")
}
