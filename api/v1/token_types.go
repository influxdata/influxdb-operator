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

// TokenSpec defines the desired state of Token
type TokenSpec struct {
	InfluxDB   string `json:"influxdb,omitempty"`
	SecretName string `json:"secret,omitempty"`
}

// TokenStatus defines the observed state of Token
type TokenStatus struct {
	Exists  bool `json:"exists,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.exists",name="Exists",type=boolean
// +kubebuilder:printcolumn:JSONPath=".status.enabled",name="Enabled",type=boolean

// Token is the Schema for zthe tokens API
type Token struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TokenSpec   `json:"spec,omitempty"`
	Status TokenStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TokenList contains a list of Token
type TokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Token `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Token{}, &TokenList{})
}
