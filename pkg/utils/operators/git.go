package operators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Shanghai-Lunara/publisher/pkg/types"
	"k8s.io/klog/v2"
)

func NewGit(gitDir string, branchName string) *Git {
	envs := make(map[string]string, 0)
	envs[types.PublisherProjectDir] = gitDir
	envs[types.PublisherGitBranch] = branchName
	return &Git{
		output: make(chan<- string, 4096),
		step: &types.Step{
			Id:             0,
			Name:           "Git-Operator",
			Phase:          types.StepPending,
			Policy:         types.StepPolicyAuto,
			Available:      types.StepAvailableEnable,
			Envs:           envs,
			Output:         make([]string, 0),
			SharingData:    make(map[string]string, 0),
			SharingSetting: false,
		},
	}
}

type Git struct {
	output chan<- string
	step   *types.Step
}

func (g *Git) Step() *types.Step {
	return g.step
}

func (g *Git) Update(s *types.Step) {
	g.step = s.DeepCopy()
}

func (g *Git) Prepare() {

}

func (g *Git) Run(output chan<- string) (res []string, err error) {
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
	if out, err = g.checkout(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	if out, err = g.pull(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	if out, err = g.source(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	if out, err = g.getCommitHash(); err != nil {
		klog.V(2).Info(err)
		g.step.Phase = types.StepFailed
		return res, err
	}
	res = append(res, string(out))
	g.step.Phase = types.StepSucceeded
	return res, nil
}

func (g *Git) cd() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s", g.step.Envs[types.PublisherProjectDir])
	return DefaultExec(commands)
}

func (g *Git) branch() (res string, err error) {
	commands := fmt.Sprintf("cd %s && Git branch -a | grep '*'", g.step.Envs[types.PublisherProjectDir])
	t, err := DefaultExec(commands)
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

func (g *Git) fetchAll() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && Git fetch --all && Git fetch -p", g.step.Envs[types.PublisherProjectDir])
	return ExecWithStreamOutput(commands, g.output)
}

func (g *Git) revert() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && Git add --all && Git checkout -f && Git reset --hard", g.step.Envs[types.PublisherProjectDir])
	return ExecWithStreamOutput(commands, g.output)
}

func (g *Git) checkout() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && Git checkout -B %s --track remotes/origin/%s",
		g.step.Envs[types.PublisherProjectDir], g.step.Envs[types.PublisherGitBranch], g.step.Envs[types.PublisherGitBranch])
	klog.Info("Git checkout commands:", commands)
	return ExecWithStreamOutput(commands, g.output)
}

func (g *Git) pull() (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && Git pull", g.step.Envs[types.PublisherProjectDir])
	return ExecWithStreamOutput(commands, g.output)
}

func (g *Git) AddAll(output chan<- string) (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && Git add --all", g.step.Envs[types.PublisherProjectDir])
	return ExecWithStreamOutput(commands, output)
}

func (g *Git) Push(output chan<- string) (res []byte, err error) {
	commands := fmt.Sprintf("cd %s && Git push", g.step.Envs[types.PublisherProjectDir])
	return ExecWithStreamOutput(commands, output)
}

func (g *Git) Commit(output chan<- string, content, source, branch, hash string) (res []byte, err error) {
	commands := fmt.Sprintf(
		`cd %s && Git commit -a -m "Automatic sync %s by lunara-publisher-robot
Srouce: %s
Branch: %s
Hash: %s"`,
		g.step.Envs[types.PublisherProjectDir],
		content,
		source,
		branch,
		hash)
	klog.Info("Git Commit commands:", commands)
	return ExecWithStreamOutput(commands, output)
}

func (g *Git) source() (res []byte, err error) {
	commands := fmt.Sprintf(`cd %s && cat .git/config | grep url`, g.step.Envs[types.PublisherProjectDir])
	res, err = DefaultExec(commands)
	if err == nil {
		g.step.Envs[types.PublisherGitSource] = string(res)
	}
	klog.Info("g.step.Envs[types.PublisherGitSource]:", g.step.Envs[types.PublisherGitSource])
	return res, err
}

func (g *Git) getCommitHash() (res []byte, err error) {
	commands := fmt.Sprintf(`cd %s && git log -p -1 | grep commit`, g.step.Envs[types.PublisherProjectDir])
	res, err = DefaultExec(commands)
	if err == nil {
		g.step.Envs[types.PublisherGitCommitHash] = string(res)
	}
	klog.Info("g.step.Envs[types.PublisherGitCommitHash]:", g.step.Envs[types.PublisherGitCommitHash])
	return res, err
}
