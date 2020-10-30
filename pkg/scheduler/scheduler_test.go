package scheduler

import (
	"reflect"
	"testing"

	"github.com/nevercase/publisher/pkg/types"
)

func TestNewScheduler(t *testing.T) {
	tests := []struct {
		name string
		want *Scheduler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScheduler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScheduler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_handle(t *testing.T) {
	type fields struct {
		items map[types.Namespace]*Groups
	}
	type args struct {
		message []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				items: tt.fields.items,
			}
			gotRes, err := s.handle(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scheduler.handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Scheduler.handle() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestScheduler_handleListNamespaces(t *testing.T) {
	type fields struct {
		items map[types.Namespace]*Groups
	}
	x := &types.Result{
		Items: []string{"hamster", "helix-2", "helix-saga"},
	}
	data, _ := x.Marshal()
	tests := []struct {
		name    string
		fields  fields
		wantRes []byte
		wantErr bool
	}{
		{
			name:    "TestScheduler_handleListNamespaces_1",
			wantRes: data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScheduler()
			gotRes, err := s.handleListNamespaces()
			if (err != nil) != tt.wantErr {
				t.Errorf("Scheduler.handleListNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Scheduler.handleListNamespaces() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestScheduler_handleListGroupNames(t *testing.T) {
	type fields struct {
		items map[types.Namespace]*Groups
	}
	type args struct {
		namespace types.Namespace
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				items: tt.fields.items,
			}
			gotRes, err := s.handleListGroupNames(tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scheduler.handleListGroupNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Scheduler.handleListGroupNames() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
