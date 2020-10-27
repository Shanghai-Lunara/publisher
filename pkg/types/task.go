package types

type Task struct {
	Id          int32   `json:"id" protobuf:"varint,1,opt,name=id"`
	Project     Project `json:"project" protobuf:"bytes,2,opt,name=project"`
	ClientSteps []Step  `json:"clientSteps" protobuf:"bytes,3,opt,name=clientSteps"`
	ServerSteps []Step  `json:"serverSteps" protobuf:"bytes,4,opt,name=serverSteps"`
}
