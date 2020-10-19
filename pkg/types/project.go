package types

type Project struct {
	Name      string            `json:"name" protobuf:"bytes,1,opt,name=name"`
	Namespace string            `json:"namespace" protobuf:"bytes,2,opt,name=namespace"`
	GroupName string            `json:"groupName" protobuf:"bytes,3,opt,name=groupName"`
	Envs      map[string]string `json:"envs" protobuf:"bytes,4,opt,name=envs"`
}
