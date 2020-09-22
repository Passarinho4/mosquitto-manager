package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Acl struct {
	AclType    string `json:"aclType" bson:"aclType"`
	AccessType string `json:"accessType" bson:"accessType"`
	Topic      string `json:"topic" bson:"topic"`
}

type MosquittoCredSpec struct {
	Id       string `json:"id" bson:"id"`
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Acls     []Acl  `json:"acls" bson:"acls"`
}

type MosquittoCred struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MosquittoCredSpec `json:"spec"`
}

type MosquittoCredList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MosquittoCred `json:"items"`
}
