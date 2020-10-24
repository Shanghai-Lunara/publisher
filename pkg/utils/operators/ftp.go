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

const (
	FtpMkdirMark = "mark"
)

type ftp struct {
	config   *operator.FtpConfig
	operator operator.FtpOperator
	step     *types.Step
}

func (f *ftp) Step() *types.Step {
	return f.step
}

func (f *ftp) Run() (res []string, err error) {
	prefix := ""
	if mark, ok := f.step.Envs[types.PublisherFtpMkdir]; ok {
		if mark == FtpMkdirMark {
			dir, err := f.yunLuoMkdir()
			if err != nil {
				klog.V(2).Info(err)
				return res, err
			}
			f.step.Envs[types.PublisherFtpMkdir] = dir
			prefix = dir
			c, err := f.operator.Conn()
			if err != nil {
				klog.V(2).Info(err)
				return res, err
			}
			if err := c.MakeDir(fmt.Sprintf("%s/%s", f.config.WorkDir, dir)); err != nil {
				klog.V(2).Info(err)
				return res, err
			}
		}
	}
	for _, v := range f.step.UploadFiles {
		target := v.TargetFile
		if prefix != "" {
			target = fmt.Sprintf("%s/%s", prefix, target)
		}
		if err := f.operator.UploadFile(v.SourceFile, target); err != nil {
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
