package scheduler

import "github.com/nevercase/publisher/pkg/types"

type Scheduler struct {
	items map[types.Namespace]Groups
}

type Groups struct {
	items map[types.GroupName]Group
}

type Group struct {
	Runners map[string]types.RunnerInfo `json:"runners" protobuf:"bytes,1,opt,name=runners"`
	Tasks   []types.Task                `json:"tasks" protobuf:"bytes,2,opt,name=tasks"`
}

func (s *Scheduler) handleListNamespaces() []string {
	return []string{"helix-saga", "helix-2", "hamster"}
}

func (s *Scheduler) handleListGroups() []string {
	return []string{"helix-saga", "helix-2", "hamster"}
}
