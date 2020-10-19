package types

type Projects struct {
	Projects map[string]Project `json:"projects" protobuf:"bytes,1,opt,name=projects"`
}

type Project struct {
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}
