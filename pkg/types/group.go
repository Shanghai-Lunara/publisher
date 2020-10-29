package types

type Namespace string

type GroupName string

type Scheduler struct {
	Items map[Namespace]Groups `json:"items" protobuf:"bytes,1,rep,name=nitems"`
}

type Groups struct {
	Items map[GroupName]Group `json:"items" protobuf:"bytes,1,rep,name=gitems"`
}

type Group struct {
	Tasks []Task `json:"tasks" protobuf:"bytes,1,rep,name=tasks"`
}
