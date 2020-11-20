package operators

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Shanghai-Lunara/go-gpt/pkg/operator"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog/v2"
)

func NewFtp(host string, port int, username, password, workDir string, timeout int) *ftp {
	envs := make(map[string]string, 0)
	envs[types.PublisherFtpHost] = host
	envs[types.PublisherFtpPort] = fmt.Sprintf("%d", port)
	envs[types.PublisherFtpUsername] = username
	envs[types.PublisherFtpPassword] = password
	envs[types.PublisherFtpWorkDir] = workDir
	envs[types.PublisherFtpTimeout] = fmt.Sprintf("%d", timeout)
	return &ftp{
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

// ftp implements github.com/nevercase/publisher/pkg/interfaces.StepOperator
type ftp struct {
	config      *operator.FtpConfig
	operator    operator.FtpOperator
	step        *types.Step
	prepareFunc func()
}

func (f *ftp) Step() *types.Step {
	return f.step
}

func (f *ftp) Update(s *types.Step) {
	f.step = s.DeepCopy()
}

func (f *ftp) SettingPrepareFunc(fc func()) {
	f.prepareFunc = fc
}

func (f *ftp) Prepare() {
	f.prepareFunc()
}

func (f *ftp) Run(output chan<- string) (res []string, err error) {
	f.step.Phase = types.StepRunning
	if err = f.ReloadConfig(); err != nil {
		klog.V(2).Info(err)
		f.step.Phase = types.StepFailed
		return nil, err
	}
	prefix := ""
	if mark, ok := f.step.Envs[types.PublisherFtpMkdir]; ok {
		if mark == FtpMkdirMark {
			dir, err := f.yunLuoMkdir()
			if err != nil {
				klog.V(2).Info(err)
				f.step.Phase = types.StepFailed
				return res, err
			}
			f.step.Envs[types.PublisherFtpMkdir] = dir
			prefix = dir
			c, err := f.operator.Conn()
			if err != nil {
				klog.V(2).Info(err)
				f.step.Phase = types.StepFailed
				return res, err
			}
			if err := c.MakeDir(fmt.Sprintf("%s/%s", f.config.WorkDir, dir)); err != nil {
				klog.V(2).Info(err)
				f.step.Phase = types.StepFailed
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
			f.step.Phase = types.StepFailed
			return res, nil
		}
	}
	f.step.Phase = types.StepSucceeded
	return res, nil
}

func (f *ftp) ReloadConfig() (err error) {
	var port, timeout int
	if port, err = strconv.Atoi(f.step.Envs[types.PublisherFtpPort]); err != nil {
		klog.V(2).Info(err)
		return err
	}
	if timeout, err = strconv.Atoi(f.step.Envs[types.PublisherFtpTimeout]); err != nil {
		klog.V(2).Info(err)
		return err
	}
	fc := &operator.FtpConfig{
		Username: f.step.Envs[types.PublisherFtpUsername],
		Password: f.step.Envs[types.PublisherFtpPassword],
		Host:     f.step.Envs[types.PublisherFtpHost],
		Port:     port,
		WorkDir:  f.step.Envs[types.PublisherFtpWorkDir],
		Timeout:  timeout,
	}
	f.config = fc
	f.operator = operator.NewFtpOperator(*fc)
	return nil
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

func (f *ftp) Operator() operator.FtpOperator {
	return f.operator
}