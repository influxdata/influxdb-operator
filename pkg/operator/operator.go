package operator

import (
	"log"
	"sync"

	"github.com/gianarb/influxdb-operator/pkg/client/tick/v1alpha1"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	VERSION = "0.0.0.dev"
)

type Options struct {
	KubeConfig string
}

type InfluxDBOperator struct {
	Options
	kubeClient *kubernetes.Clientset
	tickCs     v1alpha1.TickV1alpha1Client
	clientSet  clientset.Interface
	tickInf    cache.SharedIndexInformer
}

func (o *InfluxDBOperator) handleAddInfluxDB(obj interface{}) {
	panic("adding ugh?")
}

func New(options Options) *InfluxDBOperator {
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

	operator := &InfluxDBOperator{
		Options:    options,
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

func (operator *InfluxDBOperator) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
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
