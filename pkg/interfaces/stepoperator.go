package interfaces

import "github.com/nevercase/publisher/pkg/types"

type StepOperator interface {
	Step() *types.Step
	Update(s *types.Step)
	Run() (res []string, err error)
}
