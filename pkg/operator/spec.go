package operator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type TICKSpecList struct {
	Items []TICKSpec
}

func (t TICKSpecList) DeepCopyObject() runtime.Object {
	return t
}

func (t TICKSpecList) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}

type InfluxDBSpec struct {
	Image string `json:"image" yaml:"image"`
}

type TICKSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Type              string
	InfluxDB          InfluxDBSpec `json:"spec" yaml:"spec"`
}

func (t TICKSpec) DeepCopyObject() runtime.Object {
	return t
}

func (t TICKSpec) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}
