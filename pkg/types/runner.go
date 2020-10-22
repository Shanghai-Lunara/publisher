package types

type Runner struct {
	Hostname  string `json:"hostname" protobuf:"bytes,1,opt,name=hostname"`
	Namespace string `json:"namespace" protobuf:"bytes,2,opt,name=namespace"`
	GroupName string `json:"groupName" protobuf:"bytes,3,opt,name=groupName"`
}
