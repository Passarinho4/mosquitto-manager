package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

func (in *MosquittoCred) DeepCopyInto(out *MosquittoCred) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = MosquittoCredSpec{
		Login:    in.Spec.Login,
		Password: in.Spec.Password,
	}
	if in.Spec.Acls != nil {
		out.Spec.Acls = make([]Acl, len(in.Spec.Acls))
		for i := range in.Spec.Acls {
			in.Spec.Acls[i].DeepCopyInto(&out.Spec.Acls[i])
		}
	}
}

func (in *Acl) DeepCopyInto(out *Acl) {
	out.AccessType = in.AccessType
	out.AclType = in.AclType
	out.Topic = in.Topic
}

func (in *MosquittoCred) DeepCopyObject() runtime.Object {
	out := MosquittoCred{}
	in.DeepCopyInto(&out)
	return &out
}

func (in *MosquittoCredList) DeepCopyObject() runtime.Object {
	out := MosquittoCredList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]MosquittoCred, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}
