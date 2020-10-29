package types

type Task struct {
	Id      int32                 `json:"id" protobuf:"varint,1,opt,name=id"`
	Runners map[string]RunnerInfo `json:"runners" protobuf:"bytes,2,opt,name=runners"`
}
