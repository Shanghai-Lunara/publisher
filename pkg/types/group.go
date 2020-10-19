package types

import "sync"

type Group struct {
	mu        sync.RWMutex
	Name      string             `json:"name" protobuf:"bytes,1,opt,name=name"`
	Namespace string             `json:"namespace" protobuf:"bytes,2,opt,name=namespace"`
	Projects  map[string]Project `json:"projects" protobuf:"bytes,3,opt,name=projects"`
}
