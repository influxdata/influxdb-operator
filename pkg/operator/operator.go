package operator

import (
	"log"
	"sync"

	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
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
	config             Config
	kubeClient         *kubernetes.Clientset
	tickCs             v1alpha1.TickV1alpha1Client
	clientSet          clientset.Interface
	influxInformer     cache.SharedIndexInformer
	kapacitorInformer  cache.SharedIndexInformer
	chronografInformer cache.SharedIndexInformer
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

	registerInfluxInformer(operator)
	registerKapacitorInformer(operator)
	registerChronografInformer(operator)
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

	kapacitorCrd := extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "kapacitors" + "." + "gianarb.com",
			Labels: map[string]string{},
		},
		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   "gianarb.com",
			Version: "v1alpha1",
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: "kapacitors",
				Kind:   "kapacitor",
			},
		},
	}

	chronografCrd := extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "chronografs" + "." + "gianarb.com",
			Labels: map[string]string{},
		},
		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   "gianarb.com",
			Version: "v1alpha1",
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: "chronografs",
				Kind:   "chronograf",
			},
		},
	}

	crds := []*extensionsobj.CustomResourceDefinition{
		&influxCrd,
		&kapacitorCrd,
		&chronografCrd,
	}
	crdClient := operator.clientSet.ApiextensionsV1beta1().CustomResourceDefinitions()
	for _, crd := range crds {
		if _, err := crdClient.Create(crd); err != nil && !apierrors.IsAlreadyExists(err) {
			panic(err)
		}
	}

	go operator.kapacitorInformer.Run(stopCh)
	go operator.influxInformer.Run(stopCh)
	go operator.chronografInformer.Run(stopCh)
}
