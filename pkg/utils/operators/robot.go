package operators

import (
	"fmt"
	"github.com/Shanghai-Lunara/publisher/pkg/types"
)

func NewRobot(p types.StepPolicy, durationInMs int64) *robot {
	envs := make(map[string]string, 0)
	envs[types.RobotDurationInMs] = fmt.Sprintf("%d", durationInMs)
	return &robot{
		step: &types.Step{
			Id:             0,
			Name:           "Robot",
			Phase:          types.StepPending,
			Policy:         p,
			Available:      types.StepAvailableEnable,
			Envs:           envs,
			Output:         make([]string, 0),
			SharingData:    make(map[string]string, 0),
			SharingSetting: false,
		},
	}
}

// robot implements github.com/Shanghai-Lunara/publisher/pkg/interfaces.StepOperator
type robot struct {
	step        *types.Step
	prepareFunc func()
}

func (r *robot) Step() *types.Step {
	return r.step
}

func (r *robot) Update(s *types.Step) {
	r.step = s.DeepCopy()
}

func (r *robot) SettingPrepareFunc(fc func()) {
	r.prepareFunc = fc
}

func (r *robot) Prepare() {
	r.prepareFunc()
}

func (r *robot) Run(output chan<- string) (res []string, err error) {
	return res, nil
}
