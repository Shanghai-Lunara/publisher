package runner

import (
	"github.com/nevercase/publisher/pkg/interfaces"
	"github.com/nevercase/publisher/pkg/types"
)

type Runner struct {
	Name          string                    `json:"name" protobuf:"bytes,1,opt,name=name"`
	Hostname      string                    `json:"hostname" protobuf:"bytes,2,opt,name=hostname"`
	Namespace     string                    `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	GroupName     string                    `json:"groupName" protobuf:"bytes,4,opt,name=groupName"`
	StepOperators []interfaces.StepOperator `json:"stepOperators" protobuf:"bytes,5,opt,name=stepOperators"`
}

func (r *Runner) Run(s *types.Step) (err error) {

	return nil
}

func (r *Runner) Update(s *types.Step) (err error) {

	return nil
}
