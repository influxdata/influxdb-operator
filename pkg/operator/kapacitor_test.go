package operator

import (
	"testing"

	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
)

func TestKapacitorService(t *testing.T) {
	config := Config{
		Labels: map[string]string{
			"operator": "influxdb-operator",
		},
	}
	service := makeKapacitorService("hello-kapa", config)
	if service.GetName() != "hello-kapa" {
		t.Errorf("Expected hello-kapa as service name instead of %s.", service.GetName())
	}
	if service.Spec.Ports[0].Name != "kapacitor" {
		t.Errorf("Expected service port name kapacitor instead of %s.", service.Spec.Ports[1].Name)
	}
	if service.Spec.Ports[0].Port != 9092 {
		t.Errorf("Expected service port 9092 instead of %d.", service.Spec.Ports[1].Port)
	}

	if service.Spec.Selector["name"] != "hello-kapa" {
		t.Errorf("Expected selector name for hello-kapa instead of %s.", service.Spec.Selector["name"])
	}
}

func TestKapacitorDeployment(t *testing.T) {
	spec := v1alpha1.Kapacitor{
		Spec: v1alpha1.KapacitorSpec{
			BaseImage: "docker.io/library/kapacitor:1.2.4",
		},
	}
	deployment := makeKapacitorDeployment("hello-kapa", &spec)
	if deployment.Name != "hello-kapa" {
		t.Errorf("Expcted deployment name hello-kapa instead of %s.", deployment.Name)
	}
	if deployment.Spec.Template.Spec.Containers[0].Image != "docker.io/library/kapacitor:1.2.4" {
		t.Errorf("Expcted deployment image docker.io/library/kapacitor:1.2.4 instead of %s.", deployment.Spec.Template.Spec.Containers[0].Image)
	}

	if deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort != 9092 {
		t.Error("Expcted container port 9092 open but it's not.")
	}
}
