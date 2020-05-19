package internal

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	v1alpha12 "mosquitto-manager/internal/api/types/v1alpha1"
	"mosquitto-manager/internal/clientset/v1alpha1"
)

type ClientManager struct {
	client *v1alpha1.ExampleV1Alpha1Client
}

func NewClientManager(kubeconfig *string) *ClientManager {
	result := ClientManager{}
	_ = v1alpha12.AddToScheme(scheme.Scheme)
	client, err := v1alpha1.NewForConfig(createConfig(kubeconfig))
	if err != nil {
		panic(err.Error())
	}
	result.client = client
	return &result
}

func removeCRD(lp LoginPassword) {

}

func createConfig(kubeconfig *string) *rest.Config {
	if *kubeconfig == "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		return config
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		return config
	}
}

func (client *ClientManager) getMosquittoCreds() []LoginPassword {
	pods, _ := client.client.MosquittoCreds("default").List(metav1.ListOptions{})
	var result []LoginPassword
	for _, item := range pods.Items {
		result = append(result, LoginPassword{Login: item.Spec.Login, Password: item.Spec.Password})
	}
	return result
}

func (client *ClientManager) createMosquittoCred(lp LoginPassword) error {
	creds := v1alpha12.MosquittoCred{
	TypeMeta:metav1.TypeMeta{
		Kind: "MosquittoCred",
	},
	ObjectMeta:metav1.ObjectMeta{
		Name: lp.Login,
	},
	Spec:v1alpha12.MosquittoCredSpec{
		Login: lp.Login,
		Password: lp.Password,
		},
	}
	_, err := client.client.MosquittoCreds("default").Create(&creds)
	return err
}
