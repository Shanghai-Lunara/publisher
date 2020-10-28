package interfaces

import "github.com/nevercase/publisher/pkg/types"

type StepOperator interface {
	Step() *types.Step
	Update(s *types.Step)
	Run(output chan<- string) (res []string, err error)
}
