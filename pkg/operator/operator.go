package operator

import (
	"fmt"
	"log"
	"sync"

	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	KubeConfig string
	Labels     map[string]string
}

type Operator struct {
	config     Config
	kubeClient *kubernetes.Clientset
	tickCs     v1alpha1.TickV1alpha1Client
	clientSet  clientset.Interface
	tickInf    cache.SharedIndexInformer
}

func (o *Operator) getObject(obj interface{}) (metav1.Object, bool) {
	ts, ok := obj.(cache.DeletedFinalStateUnknown)
	if ok {
		obj = ts.Obj
	}

	oret, err := meta.Accessor(obj)
	if err != nil {
		log.Print(err)
		return nil, false
	}
	return oret, true
}

func (o *Operator) handleDeleteInfluxDB(obj interface{}) {
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	orphanDependents := false
	deploymentName := fmt.Sprintf("%s-%s", v1alpha1.InfluxDBPlural, oret.GetName())
	err := o.kubeClient.ExtensionsV1beta1().Deployments(oret.GetNamespace()).Delete(deploymentName, &metav1.DeleteOptions{
		OrphanDependents: &orphanDependents, //TODO(fntlnz): orphan dependents is now deprecated, support the new way too!
	})
	if err != nil {
		log.Printf("Error deleting deployment %s. %s", deploymentName, err)
	}
}

func (o *Operator) handleAddInfluxDB(obj interface{}) {
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	influxdbSpec := obj.(*v1alpha1.Influxdb)

	labels := o.config.Labels
	labels["name"] = fmt.Sprintf("%s-%s", v1alpha1.InfluxDBPlural, oret.GetName())
	labels["resource"] = v1alpha1.InfluxDBPlural

	replicas := int32(1)
	deployment := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      labels["name"],
			Labels:    labels,
			Namespace: oret.GetNamespace(),
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    labels,
					Namespace: oret.GetNamespace(),
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:            oret.GetName(),
							Image:           influxdbSpec.Spec.BaseImage,
							ImagePullPolicy: "Always",
							Env:             []v1.EnvVar{},
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									Name:          "http",
									ContainerPort: 8086,
									Protocol:      v1.ProtocolTCP,
								},
							},
							//VolumeMounts: []v1.VolumeMount{},
							//Resources:    v1.ResourceRequirements{},
						},
					},
					//Volumes: []v1.Volume{},
				},
			},
		},
	}

	_, err := o.kubeClient.AppsV1beta1().Deployments(oret.GetNamespace()).Create(deployment)
	fmt.Print(err) //TODO: handle in a different way?
}

func New(options Config) *Operator {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	if options.KubeConfig != "" {
		rules.ExplicitPath = options.KubeConfig
	}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		log.Fatalf("Couldn't get Kubernetes default config: %s", err)
	}

	kubeClient := kubernetes.NewForConfigOrDie(config)
	cs := clientset.NewForConfigOrDie(config)

	rest, err := v1alpha1.NewForConfig(config)
	if err != nil {
		log.Fatalf("Couldn't get Tick trd client: %s", err)
	}

	operator := &Operator{
		config:     options,
		tickCs:     *rest,
		kubeClient: kubeClient,
		clientSet:  cs,
	}

	operator.tickInf = cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc:  operator.tickCs.InfluxDBs(metav1.NamespaceAll).List,
			WatchFunc: operator.tickCs.InfluxDBs(metav1.NamespaceAll).Watch,
		},
		&v1alpha1.Influxdb{}, 0, cache.Indexers{},
	)

	operator.tickInf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    operator.handleAddInfluxDB,
		UpdateFunc: nil,
		DeleteFunc: operator.handleDeleteInfluxDB,
	})

	return operator
}

func (operator *Operator) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	influxCrd := extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "influxdbs" + "." + "gianarb.com",
			Labels: map[string]string{},
		},
		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   "gianarb.com",
			Version: "v1alpha1",
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: "influxdbs",
				Kind:   "influxdb",
			},
		},
	}

	crds := []*extensionsobj.CustomResourceDefinition{
		&influxCrd,
	}
	crdClient := operator.clientSet.ApiextensionsV1beta1().CustomResourceDefinitions()
	for _, crd := range crds {
		if _, err := crdClient.Create(crd); err != nil && !apierrors.IsAlreadyExists(err) {
			panic(err)
		}
	}

	go operator.tickInf.Run(stopCh)
}
