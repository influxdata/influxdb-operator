package k8sutil

import (
	"log"
	"time"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sWrap struct {
	kubeClient kubernetes.Interface
	crdClient  apiextensionsclient.Interface
}

func NewK8sWrap(client kubernetes.Interface, crdClient apiextensionsclient.Interface) K8sWrap {
	return K8sWrap{
		kubeClient: client,
		crdClient:  crdClient,
	}
}

func (k *K8sWrap) RegisterThridPartyResource(groupName, name, version string) error {
	_, err := k.crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			crdObject := &apiextensionsv1beta1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: name + "." + groupName},
				Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
					Group:   groupName,
					Version: version,
					Scope:   apiextensionsv1beta1.NamespaceScoped,
					Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
						Plural: name,
						Kind:   name,
					},
				},
			}

			_, err := k.crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crdObject)
			if err != nil {
				panic(err)
			}

			err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
				createdCRD, err := k.crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})
				if err != nil {
					return false, err
				}
				for _, cond := range createdCRD.Status.Conditions {
					switch cond.Type {
					case apiextensionsv1beta1.Established:
						if cond.Status == apiextensionsv1beta1.ConditionTrue {
							return true, err
						}
					case apiextensionsv1beta1.NamesAccepted:
						if cond.Status == apiextensionsv1beta1.ConditionFalse {
							log.Printf("Name conflict: %v\n", cond.Reason)
						}
					}
				}
				return false, err
			})

			if err != nil {
				deleteErr := k.crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(name, nil)
				if deleteErr != nil {
					return errors.NewAggregate([]error{err, deleteErr})
				}
				return err
			}
		} else {
			panic(err)
		}
	}
	return nil
}

func NewThridPartyResourceClient(cfg *rest.Config, groupName, version string, spec runtime.Object, specList runtime.Object) (*rest.RESTClient, *runtime.Scheme, error) {
	gv := schema.GroupVersion{Group: groupName, Version: version}
	scheme := runtime.NewScheme()
	typeBuilder := runtime.NewSchemeBuilder(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(
			gv,
			spec,
			specList,
		)
		metav1.AddToGroupVersion(scheme, gv)
		return nil
	})
	if err := typeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	config := *cfg
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}

	return client, scheme, nil
}
