package types

type Runner struct {
	Name      string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Hostname  string `json:"hostname" protobuf:"bytes,2,opt,name=hostname"`
	Namespace string `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	GroupName string `json:"groupName" protobuf:"bytes,4,opt,name=groupName"`
}
