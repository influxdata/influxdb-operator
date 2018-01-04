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

	// The IP that chronograf listens on (default: 0.0.0.0)
	Host string `json:"host"`

	// The port that chronograf listens on for insecure connections (default: 8888).
	Port int32 `json:"port"`

	// The file path to PEM encoded public key certificate.
	Cert string `json:"cert"`

	// Run the chronograf server in develop mode.
	Develop bool `json:"develop"`

	// The file path to the boltDB file (default: /var/lib/chronograf/chronograf-v1-.db).
	BoltPath string `json:"bolt_path"`

	// The file path to private key associated with given certificate.
	Key             string           `json:"key"`
	InfuxDBSource   ChronografSource `json:"influxdb_source"`
	KapacitorSource ChronografSource `json:"kapacitor_source"`

	// The path to the directory for pre-created dashboards (default: /usr/share/chronograf/canned).
	CannedPath string `json:"canned_path"`

	// The secret for signing tokens.
	TokenSecret string `json:"token_secret"`

	// The total duration (in hours) of cookie life for authentication (default:
	// 720h). Authentication expires on browser close if --auth-duration is set to
	// 0.
	AuthDuration string `json:"auth_duration"`

	// Set the logging level (default: info, accepted: debug|info|error)
	LogLevel string `json:"log_level"`
}

type ChronografSource struct {
	// The location of your InfluxDB instance including http://, the IP address, and port.
	// Example: http:///0.0.0.0:8086.
	Url string `json:"url"`
	// The username for your InfluxDB instance.
	Username string `json:"username"`
	// The password for your InfluxDB instance.
	Password string `json:"password"`
}

func (i *ChronografSpec) GetObjectKind() schema.ObjectKind {
	panic("not implemented")
}

func (i *ChronografSpec) DeepCopyObject() runtime.Object {
	panic("not implemented")
}
