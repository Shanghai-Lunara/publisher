package operators

import (
	"github.com/nevercase/publisher/pkg/interfaces"
	"github.com/nevercase/publisher/pkg/types"
)

func NewGit(gitDir string, branchName string) interfaces.StepOperator {
	envs := make(map[string]string, 0)
	envs[types.PublisherProjectDir] = gitDir
	envs[types.PublisherGitBranch] = branchName
	return &git{
		step: &types.Step{
			Id:     0,
			Name:   "Git-Operator",
			Phase:  types.StepPending,
			Policy: types.StepPolicyAuto,
			Envs:   envs,
			Output: make([]string, 0),
		},
	}
}

type git struct {
	step *types.Step
}

func (g *git) Step() *types.Step {
	return g.step
}
