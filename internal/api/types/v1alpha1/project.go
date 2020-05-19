package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type MosquittoCredSpec struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type MosquittoCred struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec MosquittoCredSpec `json:"spec"`
}

type MosquittoCredList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items []MosquittoCred `json:"items"`
}
