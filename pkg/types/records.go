package types

type Record struct {
	Id         int32     `json:"id" protobuf:"varint,1,opt,name=id"`
	Namespace  Namespace `json:"namespace" protobuf:"bytes,2,opt,name=namespace"`
	GroupName  GroupName `json:"groupName" protobuf:"bytes,3,opt,name=groupName"`
	RunnerName string    `json:"runnerName" protobuf:"bytes,4,opt,name=runnerName"`
	StepInfo   []byte    `json:"stepInfo" protobuf:"bytes,5,opt,name=stepInfo"`
	CreatedTM  int32     `json:"createdTM" protobuf:"varint,6,opt,name=createdTM"`
}
