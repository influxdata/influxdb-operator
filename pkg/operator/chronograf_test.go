package operator

import (
	"testing"

	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
)

func TestDefaultValuesForChronografSpec(t *testing.T) {
	spec := &v1alpha1.Chronograf{}
	setDefaultSpecValues(spec)
	if spec.Spec.Host != "0.0.0.0" {
		t.Errorf("Expected value 0.0.0.0 istead of %s.", spec.Spec.Host)
	}

	if spec.Spec.Port != 8888 {
		t.Errorf("Expected value 8888 istead of %d.", spec.Spec.Port)
	}
}

func TestDefaultValuesForChronografSpecPrePopulated(t *testing.T) {
	spec := &v1alpha1.Chronograf{}
	spec.Spec.Port = 9088
	setDefaultSpecValues(spec)
	if spec.Spec.Host != "0.0.0.0" {
		t.Errorf("Expected value 0.0.0.0 istead of %s.", spec.Spec.Host)
	}

	if spec.Spec.Port != 9088 {
		t.Errorf("Expected value 9088 istead of %d.", spec.Spec.Port)
	}
}

func TestChronografServiceFromSpec(t *testing.T) {
	spec := &v1alpha1.Chronograf{}
	spec.Spec.Port = 9088
	config := Config{
		Labels: map[string]string{
			"foo": "bar",
		},
	}
	service := makeChronografService("hello", config, spec)
	if len(service.GetLabels()) != 1 {
		t.Errorf("Expected 1 label for this service. We get %d : %v.", len(service.GetLabels()), service.GetLabels())
	}

	if service.GetName() != "hello" {
		t.Errorf("Expected hello as service name instead of %s.", service.GetName())
	}
	if service.Spec.Ports[0].Name != "ui" {
		t.Errorf("Expected ui as service name instead of %s.", service.Spec.Ports[0].Name)
	}

	if service.Spec.Ports[0].Port != 9088 {
		t.Errorf("Expected 9088 as service port instead of %d.", service.Spec.Ports[0].Port)
	}
}
