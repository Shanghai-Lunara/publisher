package scheduler

import (
	"fmt"
	"sort"
	"sync"

	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog"
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
	mu    sync.Mutex
	items map[types.Namespace]*Groups
}

type Groups struct {
	items map[types.GroupName]*Group
}

type Group struct {
	Runners map[string]*types.RunnerInfo `json:"runners" protobuf:"bytes,1,opt,name=runners"`
	Tasks   []*types.Task                `json:"tasks" protobuf:"bytes,2,opt,name=tasks"`
}

func (s *Scheduler) handle(message []byte) (res []byte, err error) {
	req := &types.Request{}
	if err = req.Unmarshal(message); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	switch req.Type.ServiceAPI {
	case types.ListNamespace:
		res, err = s.handleListNamespaces(req.Data)
	case types.ListGroupName:
		res, err = s.handleListGroupNames(req.Data)
	case types.ListTask:
		res, err = s.handleListTasks(req.Data)
	}
	if err != nil {
		klog.V(2).Info(err)
		//todo handle error
		return nil, err
	}
	result := &types.Response{
		Code:    0,
		Message: "",
		Type:    req.Type,
		Data:    res,
	}
	return result.Marshal()
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

func (s *Scheduler) handleRegisterRunner(data []byte) (res []byte, err error) {
	req := &types.RegisterRunnerRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	var g *Group
	if g, err = s.getGroup(req.RunnerInfo.Namespace, req.RunnerInfo.GroupName); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	g.Runners[req.RunnerInfo.Name] = &req.RunnerInfo
	result := &types.RegisterRunnerResponse{}
	return result.Marshal()
}

const (
	ErrNamespaceWasNotExisted = "error: namespace:%s was not existed"
	ErrGroupNameWasNotExisted = "error: namespace:%s groupName:%s was not existed"
)

func (s *Scheduler) getGroup(namespace types.Namespace, groupName types.GroupName) (*Group, error) {
	if t, ok := s.items[namespace]; ok {
		if t2, ok := t.items[groupName]; ok {
			return t2, nil
		} else {
			return nil, fmt.Errorf(ErrGroupNameWasNotExisted, namespace, groupName)
		}
	} else {
		return nil, fmt.Errorf(ErrNamespaceWasNotExisted, namespace)
	}
}