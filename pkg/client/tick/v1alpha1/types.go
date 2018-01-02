package v1alpha1

import (
	"k8s.io/api/core/v1"
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
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,14,opt,name=imagePullPolicy,casttype=PullPolicy"`
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this deployment.
	// It must match the pod template's labels.
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,2,opt,name=selector"`
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
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,14,opt,name=imagePullPolicy,casttype=PullPolicy"`
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this deployment.
	// It must match the pod template's labels.
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,2,opt,name=selector"`
}

func (i *KapacitorSpec) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *KapacitorSpec) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type Chronograf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ChronografSpec `json:"spec"`
}

func (i *Chronograf) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *Chronograf) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type ChronografList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*Chronograf `json:"items"`
}

func (i *ChronografList) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *ChronografList) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type ChronografSpec struct {
	BaseImage string `json:"baseImage,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,14,opt,name=imagePullPolicy,casttype=PullPolicy"`
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this deployment.
	// It must match the pod template's labels.
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,2,opt,name=selector"`
}

func (i *ChronografSpec) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *ChronografSpec) DeepCopyObject() runtime.Object {
	panic("not implemented")
}
