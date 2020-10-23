package operators

import (
	"context"
	"fmt"
	"github.com/nevercase/publisher/pkg/interfaces"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog"
	"os/exec"
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

func (g *git) Run() (res []string, err error) {
	g.step.Phase = types.StepRunning
	var out []byte
	if out, err = g.pull(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	res = append(res, string(out))
	g.step.Phase = types.StepSucceeded
	return res, nil
}

func (g *git) pull() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && git pull", g.step.Envs[types.PublisherProjectDir])
	return g.exec(commands)
}

func (g *git) exec(commands string) (res []byte, err error) {
	return exec.CommandContext(context.Background(), "sh", "-c", commands).Output()
}
