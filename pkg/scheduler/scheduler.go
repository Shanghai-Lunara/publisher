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
	case types.ServiceAPIListVersionsRequest:
		reqType.ServiceAPI = types.ServiceAPIListVersionsResponse
		res, err = s.handleListRecordsRequest(req.Data)
	}
	if err != nil {
		klog.V(2).Info(err)
		//todo handle error
		return nil, err
	}
	// todo remove
	switch req.Type.ServiceAPI {
	case types.RegisterRunner:
	case types.UpdateStep:
		if len(res) == 0 {
			return res, nil
		}
	default:
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
	keys := make([]string, 0)
	for k := range g.Runners {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, v := range keys {
		result.Runners = append(result.Runners, *g.Runners[v])
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
	var waitStep *types.Step
	for _, v := range ri.Steps {
		if v.Name == req.Step.Name {
			exist = true
			v = *req.Step.DeepCopy()
			v.Phase = types.StepRunning
			// collecting sharing data
			if v.SharingSetting == true {
				klog.Info("trigger collectSharingData name:", v.Name)
				s.collectSharingData(g, req.RunnerName, &v)
			}
			waitStep = v.DeepCopy()
			// sync for updating
			if err = s.updateStepToDashboard(req.Namespace, req.GroupName, req.RunnerName, &v); err != nil {
				klog.V(2).Info(err)
				return nil, err
			}
		}
		if exist && v.Name != req.Step.Name {
			// if the exist was true, it would change all the steps' phases to Pending
			if v.Phase != types.StepPending {
				v.Phase = types.StepPending
				if err = s.updateStepToDashboard(req.Namespace, req.GroupName, req.RunnerName, v.DeepCopy()); err != nil {
					klog.V(2).Info(err)
					return nil, err
				}
			}
		}
		newSteps = append(newSteps, v)
	}
	if exist {
		ri.Steps = newSteps
		// run the step which waited before
		go func() {
			if err = s.runStepToRunner(req.Namespace, req.GroupName, req.RunnerName, waitStep); err != nil {
				klog.V(2).Info(err)
			}
		}()
	} else {
		return nil, fmt.Errorf(ErrStepWasNotExisted, req.Namespace, req.GroupName, req.RunnerName, req.Step.Name)
	}
	return res, nil
}

func (s *Scheduler) collectSharingData(g *Group, filterRunnerName string, step *types.Step) {
	klog.V(5).Info("step:", *step)
	klog.V(5).Info("step SharingData:", step.SharingData)
	if len(step.SharingData) == 0 {
		step.SharingData = make(map[string]string, 0)
	}
	for _, v := range g.Runners {
		if v.Name == filterRunnerName {
			continue
		}
		for _, v2 := range v.Steps {
			for k, v3 := range v2.SharingData {
				step.SharingData[k] = v3
			}
		}
	}
	klog.V(5).Info("collectSharingData:", step.SharingData)
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
			// check Step Policy for automatic running when the body was types.BodyRunner
			if v.Available != types.StepAvailableDisable {
				next = false
				if v.Policy == types.StepPolicyAuto {
					// trigger running
					tn.next = true
					tn.ri = ri
					tn.step = v.DeepCopy()
					tn.step.RunnerName = req.RunnerName
					klog.V(3).Info("+++++ auto trigger step:", v.Name)
				}
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
			// send to the Runner, and then sync to all dashboards for updating Runner status
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
	req := &types.RunStepRequest{
		Namespace:  ri.Namespace,
		GroupName:  ri.GroupName,
		RunnerName: step.RunnerName,
		Step:       *step,
	}
	data, err := req.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return s.handleRunStep(data)
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
	_, err = tx.Query("INSERT INTO records (`namespace`,`groupName`,`runnerName`,`stepInfo`,`stepType`,`createdTM`) values (?,?,?,?,?,?)",
		ri.Namespace,
		ri.GroupName,
		ri.Name,
		data,
		getStepType(step),
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

func getStepType(s *types.Step) int {
	if _, ok := s.Envs[types.VersionFlag]; !ok {
		return types.RecordDefault
	}
	return types.RecordVersion
}

func (s *Scheduler) handleListRecordsRequest(data []byte) (res []byte, err error) {
	req := &types.ListRecordsRequest{}
	if err := req.Unmarshal(data); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	db := s.dao.Mysql.Master()
	var (
		rows *sql.Rows
	)
	switch req.IsVersion {
	case types.RecordDefault:
		rows, err = db.Query("SELECT * FROM records WHERE `namespace` = ? AND `groupName` = ? ORDER BY id DESC LIMIT ?, ?",
			req.Namespace,
			req.GroupName,
			req.Page,
			req.Length)
	case types.RecordVersion:
		rows, err = db.Query("SELECT * FROM records WHERE `namespace` = ? AND `groupName` = ? AND `stepType` = ? ORDER BY id DESC LIMIT ?, ?",
			req.Namespace,
			req.GroupName,
			req.IsVersion,
			req.Page,
			req.Length)
	}
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	records := make([]types.Record, 0)
	for rows.Next() {
		record := &types.Record{}
		if err := rows.Scan(&record.Id, &record.Namespace, &record.GroupName, &record.RunnerName, &record.StepInfo, &record.StepType, &record.CreatedTM); err != nil {
			klog.V(2).Info(err)
			return nil, err
		}
		records = append(records, *record)
	}
	var num int
	switch req.IsVersion {
	case types.RecordDefault:
		if err = db.QueryRow("SELECT count(*) FROM records").Scan(&num); err != nil {
			klog.V(2).Info(err)
			return nil, err
		}
	case types.RecordVersion:
		if err = db.QueryRow("SELECT count(*) FROM records WHERE `stepType` = ?", req.IsVersion).Scan(&num); err != nil {
			klog.V(2).Info(err)
			return nil, err
		}
	}
	response := &types.ListRecordsResponse{
		Params:       *req,
		Records:      records,
		RecordNumber: int32(num),
	}
	return response.Marshal()
}
