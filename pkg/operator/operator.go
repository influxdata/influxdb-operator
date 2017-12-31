package operator

import (
	"encoding/json"
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

var (
	VERSION = "0.0.0.dev"
)

type Config struct {
	KubeConfig string
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
		//c.logger.Log("msg", "get object failed", "err", err)
		return nil, false
	}
	return oret, true
}

func (o *Operator) handleAddInfluxDB(obj interface{}) {
	oret, ok := o.getObject(obj)

	if !ok {
		//TODO: error? panic? how?
		return
	}

	//TODO: There is a costant in some place
	lastApplied := oret.GetAnnotations()["kubectl.kubernetes.io/last-applied-configuration"]

	println(oret.GetName())
	var influxdbSpec v1alpha1.Influxdb
	err := json.Unmarshal([]byte(lastApplied), &influxdbSpec)
	if err != nil {
		panic(err)
	}
	replicas := int32(1)
	deployment := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pappardella",
			Labels: map[string]string{
				"name": "pappardella",
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "pappardella",
					},
					//Annotations: map[string]string{},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:            "funghiporcini",
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

	_, err = o.kubeClient.AppsV1beta1().Deployments("default").Create(deployment)
	fmt.Print(err)
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
		DeleteFunc: nil,
	})

	return operator
}

func (operator *Operator) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	log.Printf("TICK OSS operator started. Version %v\n", VERSION)

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
		log.Printf("CRD created %s", crd.Spec.Names.Kind)
	}

	go operator.tickInf.Run(stopCh)
}
