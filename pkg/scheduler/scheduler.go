package scheduler

import (
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog"
	"sort"
)

const (
	NamespaceHelixSaga types.Namespace = "helix-saga"
	NamespaceHelix2    types.Namespace = "helix-2"
	NamespaceHamster   types.Namespace = "hamster"
)

const (
	GroupNameCNLeiTing types.GroupName = "cn-leiting"
	GroupNameTWSpade   types.GroupName = "tw-spade"
)

func NewScheduler() *Scheduler {
	s := &Scheduler{
		items: make(map[types.Namespace]*Groups, 0),
	}
	s.items[NamespaceHelixSaga] = &Groups{
		items: make(map[types.GroupName]*Group, 0),
	}
	s.items[NamespaceHelixSaga].items[GroupNameCNLeiTing] = &Group{
		Runners: make(map[string]*types.RunnerInfo, 0),
		Tasks:   make([]*types.Task, 0),
	}
	s.items[NamespaceHelix2] = &Groups{
		items: make(map[types.GroupName]*Group, 0),
	}
	s.items[NamespaceHamster] = &Groups{
		items: make(map[types.GroupName]*Group, 0),
	}
	return s
}

type Scheduler struct {
	items map[types.Namespace]*Groups
}

type Groups struct {
	items map[types.GroupName]*Group
}

type Group struct {
	Runners map[string]*types.RunnerInfo `json:"runners" protobuf:"bytes,1,opt,name=runners"`
	Tasks   []*types.Task                `json:"tasks" protobuf:"bytes,2,opt,name=tasks"`
}

const (
	RequestTypeBodyError = "error req.Type.Body:%v was not matched"
)

func (s *Scheduler) handle(message []byte) (res []byte, err error) {
	req := &types.Request{}
	if err = req.Unmarshal(message); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	switch req.Type.ServiceAPI {
	case types.ListNamespace:
		return s.handleListNamespaces(req.Data)
	case types.ListGroupName:
		return s.handleListGroupNames(req.Data)
	case types.ListTask:
		return s.handleListTasks(req.Data)
	}
	return res, nil
}

func (s *Scheduler) handleListNamespaces(data []byte) (res []byte, err error) {
	keys := make([]string, 0)
	for k := range s.items {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	result := &types.ListNamespaceResponse{
		Items: keys,
	}
	return result.Marshal()
}

func (s *Scheduler) handleListGroupNames(data []byte) (res []byte, err error) {
	req := &types.ListGroupNameRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	result := &types.ListGroupNameResponse{
		Items: make([]string, 0),
	}
	if t, ok := s.items[req.Namespace]; ok {
		keys := make([]string, 0)
		for k := range t.items {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		result.Items = keys
	}
	return result.Marshal()
}

func (s *Scheduler) handleListTasks(data []byte) (res []byte, err error) {
	req := &types.ListTaskRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	result := &types.ListTaskResponse{
		Tasks: make([]types.Task, 0),
	}

	return result.Marshal()
}
