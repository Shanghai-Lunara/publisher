package scheduler

import (
	"fmt"
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
	if _, err = req.MarshalTo(message); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	switch req.Type.Body {
	case types.BodyDashboard:
		switch req.Type.Param.Command {
		case types.CommandList:
			switch req.Type.Param.ListContent {
			case types.ListContentNamespaces:
				return s.handleListNamespaces()
			case types.ListContentGroups:
				return s.handleListGroupNames(types.Namespace(req.Type.Namespace))
			case types.ListContentTasks:
			case types.ListContentSteps:
			}
		}
	case types.BodyRunner:
		switch req.Type.Param.Command {
		case types.CommandUpdate:
		case types.CommandRun:
		case types.CommandResult:
		}
	default:
		return res, fmt.Errorf(RequestTypeBodyError, req.Type.Body)
	}
	return res, nil
}

func (s *Scheduler) handleListNamespaces() (res []byte, err error) {
	keys := make([]string, 0)
	for k := range s.items {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	result := &types.Result{
		Items: keys,
	}
	return result.Marshal()
}

func (s *Scheduler) handleListGroupNames(namespace types.Namespace) (res []byte, err error) {
	result := &types.Result{
		Items: make([]string, 0),
	}
	if t, ok := s.items[namespace]; ok {
		keys := make([]string, 0)
		for k := range t.items {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		result.Items = keys
	}
	return result.Marshal()
}
