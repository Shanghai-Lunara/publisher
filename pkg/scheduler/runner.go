package scheduler

import (
	"context"
	"github.com/Shanghai-Lunara/publisher/pkg/types"
	"sync"
)

func NewRunner(ctx context.Context, info *types.RunnerInfo) *Runner {
	sub, cancel := context.WithCancel(ctx)
	r := &Runner{
		info:   info,
		ctx:    sub,
		cancel: cancel,
	}
	return r
}

type Runner struct {
	mu     sync.RWMutex
	info   *types.RunnerInfo
	ctx    context.Context
	cancel context.CancelFunc
}

func (r *Runner) UpdateStep(s *types.Step) {

}

func (r *Runner) RunStep(s *types.Step) {

}
