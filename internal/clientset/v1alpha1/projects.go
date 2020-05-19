package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"mosquitto-manager/internal/api/types/v1alpha1"
)

type MosquittoCredInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.MosquittoCredList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.MosquittoCred, error)
	Create(cred *v1alpha1.MosquittoCred) (*v1alpha1.MosquittoCred, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, opts metav1.DeleteOptions) error
}

type mosquittoCredClient struct {
	restClient rest.Interface
	ns         string
}

const mosquittoCreds = "mosquitto-creds"

func (c *mosquittoCredClient) List(opts metav1.ListOptions) (*v1alpha1.MosquittoCredList, error) {
	result := v1alpha1.MosquittoCredList{}
	err := c.restClient.
		Get().Namespace(c.ns).Resource(mosquittoCreds).VersionedParams(&opts, scheme.ParameterCodec).Do().Into(&result)
	return &result, err
}

func (c *mosquittoCredClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.MosquittoCred, error) {
	result := v1alpha1.MosquittoCred{}
	err := c.restClient.
		Get().Namespace(c.ns).Resource(mosquittoCreds).Name(name).VersionedParams(&opts, scheme.ParameterCodec).Do().Into(&result)
	return &result, err
}

func (c *mosquittoCredClient) Create(creds *v1alpha1.MosquittoCred) (*v1alpha1.MosquittoCred, error) {
	result := v1alpha1.MosquittoCred{}
	err := c.restClient.Post().Namespace(c.ns).Resource(mosquittoCreds).Body(creds).Do().Into(&result)
	return &result, err
}

func (c *mosquittoCredClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().Namespace(c.ns).Resource(mosquittoCreds).VersionedParams(&opts, scheme.ParameterCodec).Watch()
}

func (c *mosquittoCredClient) Delete(name string, opts metav1.DeleteOptions) error {
	return c.restClient.
		Delete().Namespace(c.ns).Resource(mosquittoCreds).Name(name).Body(&opts).Do().Error()
}
