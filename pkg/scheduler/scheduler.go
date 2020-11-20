package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nevercase/publisher/pkg/dao"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog/v2"
	"sort"
	"sync"
	"time"
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

func NewScheduler(broadcast chan *broadcast, d *dao.Dao) *Scheduler {
	s := &Scheduler{
		items:     make(map[types.Namespace]*Groups, 0),
		broadcast: broadcast,
		dao:       d,
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
	mu        sync.Mutex
	dao       *dao.Dao
	items     map[types.Namespace]*Groups
	broadcast chan<- *broadcast
}

type Groups struct {
	items map[types.GroupName]*Group
}

type Group struct {
	Runners map[string]*types.RunnerInfo `json:"runners" protobuf:"bytes,1,opt,name=runners"`
	Tasks   []*types.Task                `json:"tasks" protobuf:"bytes,2,opt,name=tasks"`
}

func (s *Scheduler) handle(message []byte, clientId int32) (res []byte, err error) {
	req := &types.Request{}
	if err = req.Unmarshal(message); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	reqType := req.Type
	switch req.Type.ServiceAPI {
	case types.Ping:
		res, err = s.handlePing(req.Data, clientId)
	case types.ListNamespace:
		res, err = s.handleListNamespaces(req.Data)
	case types.ListGroupName:
		res, err = s.handleListGroupNames(req.Data)
	case types.ListTask:
		res, err = s.handleListTasks(req.Data)
	case types.ListRunner:
		res, err = s.handleListRunners(req.Data)
	case types.RegisterRunner:
		res, err = s.handleRegisterRunner(req.Data, clientId)
	case types.RunStep:
		// RunStep must be sent from the Dashboard in the Scheduler handler.
		// And then the command would be transmitted to the specific Runner.
		// At the same time, the Runner status would be changed and synced to all dashboards.
		res, err = s.handleRunStep(req.Data)
	case types.UpdateStep:
		var tn *triggerNext
		res, tn, err = s.handleUpdateStep(req.Data, req.Type.Body)
		if req.Type.Body == types.BodyRunner && tn != nil && tn.next == true {
			go func() {
				_, err := s.triggerRunStep(tn.ri, tn.step)
				if err != nil {
					klog.V(2).Info(err)
				}
			}()
		}
	case types.CompleteStep:
		// CompleteStep must be sent from the Runner in the Scheduler handler.
		res, err = s.handleCompleteStep(req.Data)
	case types.LogStream:
		// LogStream must be sent from the Runner in the Scheduler handler.
		// The output should be inserted into mysql, and sent to all dashboard at the same time
		res, err = s.handleLogStream(req.Data)
	case types.ServiceAPIListRecordsRequest:
		reqType.ServiceAPI = types.ServiceAPIListRecordsResponse
		res, err = s.handleListRecordsRequest(req.Data)
	}
	if err != nil {
		klog.V(2).Info(err)
		//todo handle error
		return nil, err
	}
	if req.Type.ServiceAPI != types.RegisterRunner && len(res) == 0 {
		return res, nil
	}
	result := &types.Request{
		Type: reqType,
		Data: res,
	}
	return result.Marshal()
}

func (s *Scheduler) handlePing(data []byte, clientId int32) (res []byte, err error) {
	s.broadcast <- &broadcast{
		bt:         broadcastTypePing,
		clientId:   clientId,
		runnerName: "",
		msg:        res,
	}
	t := &types.PongResponse{}
	return t.Marshal()
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

func (s *Scheduler) handleListRunners(data []byte) (res []byte, err error) {
	req := &types.ListRunnerRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	var g *Group
	if g, err = s.getGroup(req.Namespace, req.GroupName); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	result := &types.ListRunnerResponse{
		Runners: make([]types.RunnerInfo, 0),
	}
	for _, v := range g.Runners {
		result.Runners = append(result.Runners, *v)
	}
	return result.Marshal()
}

func (s *Scheduler) handleRegisterRunner(data []byte, clientId int32) (res []byte, err error) {
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
	if _, ok := g.Runners[req.RunnerInfo.Name]; !ok {
		g.Runners[req.RunnerInfo.Name] = &req.RunnerInfo
	}
	s.broadcast <- &broadcast{
		bt:         broadcastTypeBindRunner,
		clientId:   clientId,
		runnerName: req.RunnerInfo.Name,
		msg:        res,
	}
	result := &types.RegisterRunnerResponse{}
	return result.Marshal()
}

const (
	ErrNamespaceWasNotExisted = "error: namespace:%s was not existed"
	ErrGroupWasNotExisted     = "error: namespace:%s groupName:%s was not existed"
	ErrRunnerWasNotExisted    = "error: namespace:%s groupName:%s runner:%s was not existed"
	ErrStepWasNotExisted      = "error: namespace:%s groupName:%s runner:%s step:%s was not existed"
)

func (s *Scheduler) getGroup(namespace types.Namespace, groupName types.GroupName) (*Group, error) {
	if t, ok := s.items[namespace]; ok {
		if t2, ok := t.items[groupName]; ok {
			return t2, nil
		} else {
			return nil, fmt.Errorf(ErrGroupWasNotExisted, namespace, groupName)
		}
	} else {
		return nil, fmt.Errorf(ErrNamespaceWasNotExisted, namespace)
	}
}

func (s *Scheduler) handleRunStep(data []byte) (res []byte, err error) {
	req := &types.RunStepRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	var g *Group
	if g, err = s.getGroup(req.Namespace, req.GroupName); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	s.mu.Lock()
	var ri *types.RunnerInfo
	if t, ok := g.Runners[req.RunnerName]; !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf(ErrRunnerWasNotExisted, req.Namespace, req.GroupName, req.RunnerName)
	} else {
		s.mu.Unlock()
		ri = t
	}
	klog.Info("handleRunStep name:", req.Step.Name)
	exist := false
	newSteps := make([]types.Step, 0)
	for _, v := range ri.Steps {
		if v.Name == req.Step.Name {
			exist = true
			// todo send to the Runner, and then sync to all dashboards for updating Runner status
			v = *req.Step.DeepCopy()
			v.Phase = types.StepRunning
			// run
			if err = s.runStepToRunner(req.Namespace, req.GroupName, req.RunnerName, &v); err != nil {
				klog.V(2).Info(err)
				return nil, err
			}
			// sync for updating
			if err = s.updateStepToDashboard(req.Namespace, req.GroupName, req.RunnerName, &v); err != nil {
				klog.V(2).Info(err)
				return nil, err
			}
		}
		if exist {
			// if the exist was true, it would change all the steps' phases to Pending
			//v.Phase = types.StepPending
			//if err = s.updateStepToDashboard(req.Namespace, req.GroupName, req.RunnerName, v.DeepCopy()); err != nil {
			//	klog.V(2).Info(err)
			//	return nil, err
			//}
		}
		newSteps = append(newSteps, v)
	}
	ri.Steps = newSteps
	if !exist {
		return nil, fmt.Errorf(ErrStepWasNotExisted, req.Namespace, req.GroupName, req.RunnerName, req.Step.Name)
	}
	return res, nil
}

type triggerNext struct {
	next bool
	ri   *types.RunnerInfo
	step *types.Step
}

func (s *Scheduler) handleUpdateStep(data []byte, body types.Body) (res []byte, tn *triggerNext, err error) {
	tn = &triggerNext{
		next: false,
		ri:   &types.RunnerInfo{},
		step: &types.Step{},
	}
	req := &types.RunStepRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, tn, err
	}
	var g *Group
	if g, err = s.getGroup(req.Namespace, req.GroupName); err != nil {
		klog.V(2).Info(err)
		return nil, tn, err
	}
	s.mu.Lock()
	var ri *types.RunnerInfo
	if t, ok := g.Runners[req.RunnerName]; !ok {
		s.mu.Unlock()
		return nil, tn, fmt.Errorf(ErrRunnerWasNotExisted, req.Namespace, req.GroupName, req.RunnerName)
	} else {
		s.mu.Unlock()
		ri = t
	}
	exist := false
	next := false
	newSteps := make([]types.Step, 0)
	for _, v := range ri.Steps {
		switch next {
		case false:
			//klog.Infof("next false body:%s step-name:%s step-phase:%s current-step-name:%s", body, req.Step.Name, req.Step.Phase, v.Name)
			if v.Name == req.Step.Name {
				exist = true
				v = req.Step
				// save to db
				if body == types.BodyRunner {
					go s.recordStep(ri, v.DeepCopy())
				}
				// sync for updating
				if err = s.updateStepToDashboard(req.Namespace, req.GroupName, req.RunnerName, &v); err != nil {
					klog.V(2).Info(err)
					return nil, tn, err
				}
				// if the request body was types.BodyRunner and the step.Phase was the types.StepSucceeded,
				// it means that the Scheduler should trigger automatic running
				if body == types.BodyRunner && v.Phase == types.StepSucceeded {
					next = true
				}
			}
		case true:
			//klog.Infof("next true body:%s step-name:%s step-phase:%s current-step-name:%s", body, req.Step.Name, req.Step.Phase, v.Name)
			// todo check Step Policy for automatic running when the body was types.BodyRunner
			if v.Policy == types.StepPolicyAuto {
				if v.Phase != types.StepDisabled {
					// trigger running
					next = false
					tn.next = true
					tn.ri = ri
					tn.step = v.DeepCopy()
					klog.V(3).Info("+++++ auto trigger step:", v.Name)
				}
			} else {
				next = false
			}
		}
		newSteps = append(newSteps, v)
	}
	ri.Steps = newSteps
	if !exist {
		return nil, tn, fmt.Errorf(ErrStepWasNotExisted, req.Namespace, req.GroupName, req.RunnerName, req.Step.Name)
	}
	return res, tn, nil
}

func (s *Scheduler) handleCompleteStep(data []byte) (res []byte, err error) {
	req := &types.RunStepRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	var g *Group
	if g, err = s.getGroup(req.Namespace, req.GroupName); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	s.mu.Lock()
	var ri *types.RunnerInfo
	if t, ok := g.Runners[req.RunnerName]; !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf(ErrRunnerWasNotExisted, req.Namespace, req.GroupName, req.RunnerName)
	} else {
		s.mu.Unlock()
		ri = t
	}
	exist := false
	newSteps := make([]types.Step, 0)
	for _, v := range ri.Steps {
		if v.Name == req.Step.Name {
			exist = true
			// todo send to the Runner, and then sync to all dashboards for updating Runner status
			v = req.Step
			// sync for updating
			if err = s.updateStepToDashboard(req.Namespace, req.GroupName, req.RunnerName, &v); err != nil {
				klog.V(2).Info(err)
				return nil, err
			}
		}
		if exist {
			// todo check Step Policy for automatic running
		}
		newSteps = append(newSteps, v)
	}
	ri.Steps = newSteps
	if !exist {
		return nil, fmt.Errorf(ErrStepWasNotExisted, req.Namespace, req.GroupName, req.RunnerName, req.Step.Name)
	}
	return res, nil
}

func (s *Scheduler) runStepToRunner(namespace types.Namespace, groupName types.GroupName, runnerName string, step *types.Step) (err error) {
	req1 := &types.RunStepRequest{
		Namespace:  namespace,
		GroupName:  groupName,
		RunnerName: runnerName,
		Step:       *step,
	}
	data1, err := req1.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	req2 := &types.Request{
		Type: types.Type{
			ServiceAPI: types.RunStep,
		},
		Data: data1,
	}
	data2, err := req2.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	s.broadcast <- &broadcast{
		bt:         broadcastTypeRunner,
		runnerName: runnerName,
		msg:        data2,
	}
	return nil
}

func (s *Scheduler) updateStepToDashboard(namespace types.Namespace, groupName types.GroupName, runnerName string, step *types.Step) (err error) {
	req := &types.UpdateStepRequest{
		Namespace:  namespace,
		GroupName:  groupName,
		RunnerName: runnerName,
		Step:       *step,
	}
	data, err := req.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	req2 := &types.Request{
		Type: types.Type{
			ServiceAPI: types.UpdateStep,
		},
		Data: data,
	}
	data2, err := req2.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	s.broadcast <- &broadcast{
		bt:         broadcastTypeDashboard,
		runnerName: "",
		msg:        data2,
	}
	return nil
}

func (s *Scheduler) handleLogStream(data []byte) (res []byte, err error) {
	req := &types.LogStreamRequest{}
	if err = req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	// todo insert into the db or runtime cache

	// broadcast to all dashboards
	req2 := &types.Request{
		Type: types.Type{
			ServiceAPI: types.LogStream,
		},
		Data: data,
	}
	data2, err := req2.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	s.broadcast <- &broadcast{
		bt:         broadcastTypeDashboard,
		runnerName: "",
		msg:        data2,
	}
	return res, nil
}

func (s *Scheduler) triggerRunStep(ri *types.RunnerInfo, step *types.Step) (res []byte, err error) {
	klog.Info("triggerRunStep name:", step.Name)
	exist := false
	newSteps := make([]types.Step, 0)
	for _, v := range ri.Steps {
		if v.Name == step.Name {
			exist = true
			// todo send to the Runner, and then sync to all dashboards for updating Runner status
			v = *step.DeepCopy()
			//v.Phase = types.StepRunning
			// run
			if err = s.runStepToRunner(ri.Namespace, ri.GroupName, ri.Name, &v); err != nil {
				klog.V(2).Info(err)
				return nil, err
			}
			// sync for updating
			if err = s.updateStepToDashboard(ri.Namespace, ri.GroupName, ri.Name, &v); err != nil {
				klog.V(2).Info(err)
				return nil, err
			}
		}
		//if exist {
		//	v.Phase = types.StepPending
		//	if err = s.updateStepToDashboard(ri.Namespace, ri.GroupName, ri.Name, v.DeepCopy()); err != nil {
		//		klog.V(2).Info(err)
		//		return nil, err
		//	}
		//}
		newSteps = append(newSteps, v)
	}
	ri.Steps = newSteps
	if !exist {
		return nil, fmt.Errorf(ErrStepWasNotExisted, ri.Namespace, ri.GroupName, ri.Name, step.Name)
	}
	return res, nil
}

func (s *Scheduler) recordStep(ri *types.RunnerInfo, step *types.Step) {
	data, err := step.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return
	}
	db := s.dao.Mysql.Master()
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		klog.V(2).Info(err)
	}
	_, err = tx.Query("INSERT INTO records (`namespace`,`groupName`,`runnerName`,`stepInfo`,`createdTM`) values (?,?,?,?,?)",
		ri.Namespace,
		ri.GroupName,
		step.RunnerName,
		data,
		time.Now().Unix())
	if err != nil {
		klog.V(2).Info(err)
		if err := tx.Rollback(); err != nil {
			klog.V(2).Info(err)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		klog.V(2).Info(err)
		return
	}
}

func (s *Scheduler) handleListRecordsRequest(data []byte) (res []byte, err error) {
	req := &types.ListRecordsRequest{}
	if err := req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	db := s.dao.Mysql.Master()
	rows, err := db.Query("SELECT * FROM records WHERE `namespace` = ? AND `groupName` = ? ORDER BY id DESC LIMIT ?, ?",
		req.Namespace,
		req.GroupName,
		req.Page,
		req.Length)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	records := make([]types.Record, 0)
	for rows.Next() {
		record := &types.Record{}
		if err := rows.Scan(&record.Id, &record.Namespace, &record.GroupName, &record.RunnerName, &record.StepInfo, &record.CreatedTM); err != nil {
			klog.V(2).Info(err)
			return nil, err
		}
		records = append(records, *record)
	}
	response := &types.ListRecordsResponse{
		Params:  *req,
		Records: records,
	}
	return response.Marshal()
}
