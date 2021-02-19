package scheduler

import (
	"context"
	"fmt"
	"github.com/Shanghai-Lunara/publisher/pkg/types"
	"k8s.io/klog/v2"
	"sync/atomic"
)

func NewRunner(ctx context.Context, info *types.RunnerInfo) *Runner {
	sub, cancel := context.WithCancel(ctx)
	r := &Runner{
		info:   info,
		rwLock: 0,
		status: Idle,
		ctx:    sub,
		cancel: cancel,
	}
	return r
}

type RunnerStatus int

const (
	Running RunnerStatus = iota
	Idle    RunnerStatus = 1
)

type Runner struct {
	rwLock                    int32
	info                      *types.RunnerInfo
	status                    RunnerStatus
	ctx                       context.Context
	cancel                    context.CancelFunc
	RecordFunc                func(ri *types.RunnerInfo, step *types.Step)
	UpdateStepToDashboardFunc func(namespace types.Namespace, groupName types.GroupName, runnerName string, step *types.Step) (err error)
}

func (r *Runner) UpdateStep(req *types.Step, body types.Body) (res []byte, tn *triggerNext, err error) {
	switch body {
	case types.BodyDashboard:
		lock := atomic.AddInt32(&r.rwLock, 1)
		defer atomic.AddInt32(&r.rwLock, -1)
		if lock >= 2 {
			return
		}
		switch r.status {
		case Running:
			return
		case Idle:
		}
	case types.BodyRunner:
	}
	tn = &triggerNext{
		next: false,
		ri:   &types.RunnerInfo{},
		step: &types.Step{},
	}
	exist := false
	next := false
	newSteps := make([]types.Step, 0)
	for _, v := range r.info.Steps {
		switch next {
		case false:
			if v.Name == req.Name {
				exist = true
				v = *req
				// save to db
				if body == types.BodyRunner {
					go r.RecordFunc(r.info, v.DeepCopy())
				}
				// sync for updating
				if err = r.UpdateStepToDashboardFunc(r.info.Namespace, r.info.GroupName, req.RunnerName, &v); err != nil {
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
					tn.ri = r.info
					tn.step = v.DeepCopy()
					tn.step.RunnerName = req.RunnerName
					klog.V(3).Info("+++++ auto trigger step:", v.Name)
				}
			}
		}
		newSteps = append(newSteps, v)
	}
	r.info.Steps = newSteps
	if !exist {
		return nil, tn, fmt.Errorf(ErrStepWasNotExisted, r.info.Namespace, r.info.GroupName, req.RunnerName, req.Name)
	}
	return res, tn, nil
}

func (r *Runner) RunStep(s *types.Step) {
	lock := atomic.AddInt32(&r.rwLock, 1)
	defer atomic.AddInt32(&r.rwLock, -1)
	if lock >= 2 {
		return
	}
	r.status = Running
}
