package operator

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gianarb/influxdb-operator/pkg/k8sutil"
	"k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	VERSION    = "0.0.0.dev"
	TPRGroup   = "gianarb.com"
	TPRVersion = "v1beta1"
	TPRName    = "tick"
)

type Options struct {
	KubeConfig string
}

type InfluxDBOperator struct {
	Options
	k8swrap       k8sutil.K8sWrap
	tickTprClient *rest.RESTClient
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

	spec := TICKSpec{}
	specList := TICKSpecList{}
	tickTrpClient, _, err := k8sutil.NewThridPartyResourceClient(config, TPRGroup, TPRVersion, &spec, &specList)
	if err != nil {
		log.Fatalf("Couldn't get Tick trd client: %s", err)
	}

	return &InfluxDBOperator{
		Options:       options,
		k8swrap:       k8sutil.NewK8sWrap(kubeClient, cs),
		tickTprClient: tickTrpClient,
	}
}

func (operator *InfluxDBOperator) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	log.Printf("TICK OSS operator started. Version %v\n", VERSION)
	operator.k8swrap.RegisterThridPartyResource(TPRGroup, TPRName, TPRVersion)
	events, errChan := operator.processTickEvents(stopCh)
	go func(events <-chan *TICKSpec) {
		for {
			select {
			case event := <-events:
				fmt.Printf("RUN %+v\n", event)
			case err := <-errChan:
				fmt.Printf("RUN %s\n", err)
			case <-stopCh:
				wg.Done()
				return
			}
		}
	}(events)
}

func (operator *InfluxDBOperator) processTickEvents(stopCh <-chan struct{}) (<-chan *TICKSpec, <-chan error) {
	events := make(chan *TICKSpec)
	errc := make(chan error, 1)
	source := cache.NewListWatchFromClient(operator.tickTprClient, TPRName, v1.NamespaceAll, fields.Everything())
	createAddHandler := func(obj interface{}) {
		event := obj.(*TICKSpec)
		event.Type = "ADDED"
		events <- event
	}

	createDeleteHandler := func(obj interface{}) {
		event := obj.(*TICKSpec)
		event.Type = "DELETED"
		events <- event
	}

	updateHandler := func(old interface{}, obj interface{}) {
		event := obj.(*TICKSpec)
		event.Type = "MODIFIED"
		events <- event
	}

	_, controller := cache.NewInformer(
		source,
		&TICKSpec{},
		time.Minute*60,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    createAddHandler,
			UpdateFunc: updateHandler,
			DeleteFunc: createDeleteHandler,
		})

	go controller.Run(stopCh)

	return events, errc
}
