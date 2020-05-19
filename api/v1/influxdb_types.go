/*
Copyright 2020 InfluxData

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InfluxDBToken Config for InfluxDBSpecs
type InfluxDBToken struct {
	SecretName string `json:"secretName,omitempty"`
	SecretKey  string `json:"secretKey,omitempty"`
}

// InfluxDBSpec defines the desired state of InfluxDB
type InfluxDBSpec struct {
	// Should the operator deploy and manage this InfluxDB, or is it an existing deployment?
	Provision bool `json:"provision"`

	// The URL of the InfluxDB server (Only used when provision is false)
	// +optional
	URL string `json:"url,omitempty"`

	// The organization
	Organization string `json:"organization,omitempty"`

	// If we're not provisioning this InfluxDB through the operator, we need an admin
	// token to do further provisioning
	// +optional
	Token InfluxDBToken `json:"token,omitempty"`
}

// InfluxDBStatus defines the observed state of InfluxDB
type InfluxDBStatus struct {
	// Is this InfluxDB responding on /health?
	Available bool `json:"available"`
	// Can we authenticate against this InfluxDB?
	Authenticated bool `json:"authenticated"`

	// Stats refresh

	// +kubebuilder:default=0
	StatsRefresh int64 `json:"statsRefreshTime"`

	// DNS Resolution time
	DNSResolution int64 `json:"dnsResolutionDuration"`

	// Connect Time
	ConnectTime int64 `json:"connectTime"`

	// First Byte
	FirstByte int64 `json:"firstByte"`

	// TLS Handshake
	TLSHandshake int64 `json:"tlsHandshake"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".spec.url",name=URL,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.organization",name=Organization,type=string
// +kubebuilder:printcolumn:JSONPath=".status.available",name=Available,type=boolean
// +kubebuilder:printcolumn:JSONPath=".status.authenticated",name=Authenticated,type=boolean
// +kubebuilder:printcolumn:JSONPath=".status.dnsResolutionDuration",name="DNS Resolution (ms)",type=integer,priority=10
// +kubebuilder:printcolumn:JSONPath=".status.connectTime",name="HTTP Connect (ms)",type=integer,priority=10
// +kubebuilder:printcolumn:JSONPath=".status.tlsHandshake",name="TLS Handshake (ms)",type=integer,priority=10
// +kubebuilder:printcolumn:JSONPath=".status.firstByte",name="First Byte (ms)",type=integer,priority=10

// InfluxDB is the Schema for the influxdbs API
type InfluxDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InfluxDBSpec   `json:"spec,omitempty"`
	Status InfluxDBStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InfluxDBList contains a list of InfluxDB
type InfluxDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InfluxDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InfluxDB{}, &InfluxDBList{})
}
