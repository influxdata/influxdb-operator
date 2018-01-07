package operator

import (
	"bytes"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInfluxDBDeployment(t *testing.T) {
	replicas := int32(1)
	spec := &v1alpha1.Influxdb{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "influxdb",
				"foo": "bar",
			},
		},
		Spec: v1alpha1.InfluxdbSpec{
			BaseImage:       "docker.io/library/influxdb:1.4.0",
			ImagePullPolicy: "Always",
			Replicas:        &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{},
			},
		},
	}

	deployment := makeInfluxDBDeployment("hello", spec)
	if len(deployment.GetLabels()) != 4 {
		t.Errorf("Expected 2 labels instead of %d: %v", len(deployment.GetLabels()), deployment.GetLabels())
	}
	if deployment.GetName() != "hello" {
		t.Errorf("Expected hello as name instead of %s.", deployment.GetName())
	}
}

func TestInfluxDBConfigMap(t *testing.T) {
	spec := v1alpha1.InfluxdbSpec{}
	config := createInfluxDBConfiguration(spec)
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(&config); err != nil {
		t.Error(err)
	}
	configMap, _ := makeInfluxDBConfigMap("hello-config", config)
	if configMap.GetName() != "hello-config" {
		t.Errorf("Expected name hello-config instead of %s", configMap.GetName())
	}
	if configMap.Data["config.toml"] != buf.String() {
		t.Error("InfluxDB config is not the same injected in the ConfigMap")
	}
}

func TestInfluxDBService(t *testing.T) {
	config := Config{
		Labels: map[string]string{
			"operator": "influxdb-operator",
		},
	}
	service := makeInfluxDBService("hello-service", config)
	if service.GetName() != "hello-service" {
		t.Errorf("Expected service name hello-service instead of %s.", service.GetName())
	}
	if len(service.GetLabels()) != 1 {
		t.Errorf("Expected 1 labels instead of %d: %v.", len(service.GetLabels()), service.GetLabels())
	}
}
