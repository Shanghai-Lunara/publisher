package types

type Task struct {
	Id      int32   `json:"id" protobuf:"varint,1,opt,name=id"`
	Project Project `json:"project" protobuf:"bytes,2,opt,name=project"`
	Steps   []Step  `json:"steps" protobuf:"bytes,3,opt,name=Step"`
}
