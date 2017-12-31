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

type Kapacitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KapacitorSpec `json:"spec"`
}

func (i *Kapacitor) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *Kapacitor) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type KapacitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*Kapacitor `json:"items"`
}

func (i *KapacitorList) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *KapacitorList) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type KapacitorSpec struct {
	BaseImage string `json:"baseImage,omitempty"`
}

func (i *KapacitorSpec) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *KapacitorSpec) DeepCopyObject() runtime.Object {
	panic("not implemented")
}
