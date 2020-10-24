package operators

import (
	"fmt"
	"time"

	"github.com/Shanghai-Lunara/go-gpt/pkg/operator"
	"github.com/nevercase/publisher/pkg/interfaces"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog"
)

func NewFtp(host string, port int, username, password, workDir string, timeout int) interfaces.StepOperator {
	envs := make(map[string]string, 0)
	fc := &operator.FtpConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		WorkDir:  workDir,
		Timeout:  timeout,
	}
	return &ftp{
		config:   fc,
		operator: operator.NewFtpOperator(*fc),
		step: &types.Step{
			Id:     0,
			Name:   "Ftp-Operator",
			Phase:  types.StepPending,
			Policy: types.StepPolicyAuto,
			Envs:   envs,
			Output: make([]string, 0),
		},
	}
}

type ftp struct {
	config   *operator.FtpConfig
	operator operator.FtpOperator
	step     *types.Step
}

func (f *ftp) Step() *types.Step {
	return f.step
}

func (f *ftp) Run() (res []string, err error) {
	if _, ok := f.step.Envs[types.PublisherFtpMkdir]; ok {
		dir, err := f.yunLuoMkdir()
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		f.step.Envs[types.PublisherFtpMkdir] = dir
		c, err := f.operator.Conn()
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err := c.MakeDir(dir); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
	}
	for _, v := range f.step.UploadFiles {
		if err := f.operator.UploadFile(v.SourceFile, v.TargetFile); err != nil {
			klog.V(2).Info(err)
			return res, nil
		}
	}
	return res, nil
}

func (f *ftp) yunLuoMkdir() (dir string, err error) {
	date := time.Now().Format("20060102")
	res, err := f.operator.List(date)
	if err != nil {
		klog.V(2).Info(err)
		return dir, err
	}
	dir = fmt.Sprintf("%s_%d", date, 1+len(res))
	return dir, nil
}
