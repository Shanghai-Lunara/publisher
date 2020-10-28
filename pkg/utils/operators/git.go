package operators

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nevercase/publisher/pkg/interfaces"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog"
)

func NewGit(gitDir string, branchName string) interfaces.StepOperator {
	envs := make(map[string]string, 0)
	envs[types.PublisherProjectDir] = gitDir
	envs[types.PublisherGitBranch] = branchName
	return &git{
		output: make(chan<- string, 4096),
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
	output chan<- string
	step   *types.Step
}

func (g *git) Step() *types.Step {
	return g.step
}

func (g *git) Update(s *types.Step) {
	g.step = s.DeepCopy()
}

func (g *git) Run(output chan<- string) (res []string, err error) {
	g.output = output
	g.step.Phase = types.StepRunning
	var out []byte
	if out, err = g.cd(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	if out, err = g.revert(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	if out, err = g.pull(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	res = append(res, string(out))
	g.step.Phase = types.StepSucceeded
	return res, nil
}

func (g *git) cd() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s", g.step.Envs[types.PublisherProjectDir])
	return g.exec(commands)
}

func (g *git) branch() (res string, err error) {
	commands := fmt.Sprintf("cd %s && git branch -a | grep '*'", g.step.Envs[types.PublisherProjectDir])
	t, err := g.exec(commands)
	if err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	activeMatched, err := regexp.Match(`\*`, t)
	if err != nil {
		return res, fmt.Errorf("git regexp.Match active target-name:%s err:%v\n", g.step.Envs[types.PublisherGitBranch], err)
	}
	if activeMatched == false {
		return res, fmt.Errorf("git regexp.Match active target-name:%s failed", g.step.Envs[types.PublisherGitBranch])
	}
	var name string
	name = strings.Replace(string(t), " ", "", -1)
	name = strings.Replace(name, "*", "", -1)
	// such as `* test\n`
	name = strings.Replace(name, "\n", "", -1)
	return name, nil
}

func (g *git) fetchAll() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && git fetch --all && git fetch -p", g.step.Envs[types.PublisherProjectDir])
	return g.execWithStreamOutput(commands)
}

func (g *git) revert() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && git add --all && git checkout -f && git reset --hard", g.step.Envs[types.PublisherProjectDir])
	return g.execWithStreamOutput(commands)
}

func (g *git) checkout() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && git checkout -B %s --track remotes/origin/%s",
		g.step.Envs[types.PublisherProjectDir], g.step.Envs[types.PublisherGitBranch], g.step.Envs[types.PublisherGitBranch])
	klog.Info("git checkout commands:", commands)
	return g.execWithStreamOutput(commands)
}

func (g *git) pull() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && git pull", g.step.Envs[types.PublisherProjectDir])
	return g.execWithStreamOutput(commands)
}

func (g *git) push() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && git push", g.step.Envs[types.PublisherProjectDir])
	return g.execWithStreamOutput(commands)
}

func (g *git) exec(commands string) (res []byte, err error) {
	return exec.CommandContext(context.Background(), "sh", "-c", commands).Output()
}

func (g *git) execWithStreamOutput(commands string) (res []byte, err error) {
	//return exec.CommandContext(context.Background(), "sh", "-c", commands).Output()
	cmd := exec.CommandContext(context.Background(), "sh", "-c", commands)
	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	if err = cmd.Start(); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		g.output <- m
		fmt.Println(m)
	}
	if err = cmd.Wait(); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	return res, nil
}
