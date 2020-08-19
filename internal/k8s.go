package internal

import (
	"errors"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	v1alpha12 "mosquitto-manager/internal/api/types/v1alpha1"
	"mosquitto-manager/internal/clientset/v1alpha1"
	"time"
)

type K8sManager struct {
	client *v1alpha1.ExampleV1Alpha1Client
}

func NewK8sManager(kubeconfig *string) Manager {
	var result = K8sManager{}
	_ = v1alpha12.AddToScheme(scheme.Scheme)
	client, err := v1alpha1.NewForConfig(createConfig(kubeconfig))
	if err != nil {
		panic(err.Error())
	}
	result.client = client
	return result
}

func (client K8sManager) Remove(id Id) error {
	err := client.client.MosquittoCreds("default").Delete(id.Id, metav1.DeleteOptions{})
	return err
}

func (client K8sManager) GetAll() []CredsWithId {
	creds, _ := client.client.MosquittoCreds("default").List(metav1.ListOptions{})
	var result []CredsWithId
	for _, item := range creds.Items {
		result = append(result, CredsWithId{
			Id:       item.Spec.Id,
			Login:    item.Spec.Login,
			Password: item.Spec.Password,
			Acls:     item.Spec.Acls,
		})
	}
	return result
}

func (client K8sManager) Get(id Id) (*CredsWithId, error) {
	creds, _ := client.client.MosquittoCreds("default").Get(id.Id, metav1.GetOptions{})
	var result = CredsWithId{
		Id:       creds.Spec.Id,
		Login:    creds.Spec.Login,
		Password: creds.Spec.Password,
		Acls:     creds.Spec.Acls,
	}
	return &result, nil
}

func (client K8sManager) Update(id Id, lp Creds) error {
	return errors.New("update is not supported with CRD storage")
}

func (client K8sManager) Create(lp Creds) (*string, error) {
	id := uuid.New().String()
	creds := v1alpha12.MosquittoCred{
		TypeMeta: metav1.TypeMeta{
			Kind: "MosquittoCred",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: v1alpha12.MosquittoCredSpec{
			Id:       id,
			Login:    lp.Login,
			Password: lp.Password,
			Acls:     lp.Acls,
		},
	}
	_, err := client.client.MosquittoCreds("default").Create(&creds)
	if err != nil {
		log.Print("Error during CRD creation ", err)
		return nil, err
	} else {
		return &id, nil
	}
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

func (client K8sManager) ObserveIfSupported(service ManagerService) {
	watch, _ := client.client.MosquittoCreds("default").Watch(v1.ListOptions{})
	for range watch.ResultChan() {
		log.Printf("Received event from K8s watcher")
		service.reloadAfterChange()
	}

}

func (client K8sManager) IsObserveSupported() bool {
	return true
}

func watchMosquittoCreds(clientSet v1alpha1.ExampleV1Alpha1Interface) cache.Store {
	projectStore, projectController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.MosquittoCreds("default").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.MosquittoCreds("default").Watch(lo)
			},
		},
		&v1alpha12.MosquittoCred{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{})
	go projectController.Run(wait.NeverStop)
	return projectStore
}
