package operators

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Shanghai-Lunara/publisher/pkg/interfaces"
	"github.com/Shanghai-Lunara/publisher/pkg/types"
)

var fakeGitOperator = NewGit("/Users/nevermore/go/src/github.com/Shanghai-Lunara/client-tool", "main")

var fakeErrGitOperator = NewGit("/Users/nevermore/go/src/github.com/Shanghai-Lunara/client-tool2", "main")

var fakeGitOperator2 = NewGit("/Users/nevermore/go/src/github.com/Shanghai-Lunara/client-tool", "test")

var fakeErrGitOperator2 = NewGit("/Users/nevermore/go/src/github.com/Shanghai-Lunara/client-tool", "test2xxxxx")

func TestNewGit(t *testing.T) {
	type args struct {
		gitDir     string
		branchName string
	}
	tests := []struct {
		name string
		args args
		want interfaces.StepOperator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGit(tt.args.gitDir, tt.args.branchName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_git_Step(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name   string
		fields fields
		want   *types.Step
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			if got := g.Step(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Git.Step() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_git_Run(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []string
		wantErr bool
	}{
		{
			name:    "Test_git_Run_1",
			fields:  fields{step: fakeGitOperator.Step()},
			wantRes: nil,
			wantErr: false,
		},
		{
			name:    "Test_git_Run_2",
			fields:  fields{step: fakeErrGitOperator.Step()},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			gotRes, err := g.Run(make(chan<- string, 4096))
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range gotRes {
				fmt.Printf("Run output:%s\n", v)
			}
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("Git.Run() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}

func Test_git_pull(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			gotRes, err := g.pull()
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.pull() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Git.pull() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_git_exec(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	type args struct {
		commands string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []byte
		wantErr bool
	}{
		{
			name:    "Test_git_exec_1",
			fields:  fields{step: fakeGitOperator.Step()},
			args:    args{commands: "date"},
			wantRes: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := DefaultExec(tt.args.commands)
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("exec `%s` output:%s\n", tt.args.commands, gotRes)
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("Git.exec() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}

func Test_git_cd(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			gotRes, err := g.cd()
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.cd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Git.cd() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_git_fetchAll(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			gotRes, err := g.fetchAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.fetchAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Git.fetchAll() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_git_revert(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			gotRes, err := g.revert()
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.revert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Git.revert() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_git_checkout(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []byte
		wantErr bool
	}{
		{
			name:    "Test_git_checkout_1",
			fields:  fields{step: fakeGitOperator.Step()},
			wantRes: nil,
			wantErr: false,
		},
		{
			name:    "Test_git_checkout_2",
			fields:  fields{step: fakeErrGitOperator2.Step()},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:    "Test_git_checkout_3",
			fields:  fields{step: fakeGitOperator2.Step()},
			wantRes: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			gotRes, err := g.checkout()
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.checkout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("checkout to: ", g.step.Envs[types.PublisherGitBranch])
			for _, v := range gotRes {
				fmt.Printf("checkout out:%s", string(v))
			}
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("Git.checkout() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}

func Test_git_branch(t *testing.T) {
	type fields struct {
		step *types.Step
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes string
		wantErr bool
	}{
		{
			name:    "Test_git_branch_1",
			fields:  fields{step: fakeGitOperator.Step()},
			wantRes: "main",
			wantErr: false,
		},
		{
			name:    "Test_git_branch_2",
			fields:  fields{step: fakeGitOperator2.Step()},
			wantRes: "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{
				output: make(chan<- string, 4096),
				step:   tt.fields.step,
			}
			if _, err := g.checkout(); err != nil {
				t.Errorf("Git.branch() g.checkout() error = %v", err)
			}
			gotRes, err := g.branch()
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.branch() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotRes != tt.wantRes {
				t.Errorf("Git.branch() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
