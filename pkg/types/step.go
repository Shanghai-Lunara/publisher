package types

// StepPhase is a label for the condition of a Step at the current time.
type StepPhase string

// These are the valid statuses of Steps.
const (
	// StepPending means the Step has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	StepPending StepPhase = "Pending"
	// StepRunning means the Step has been bound to a Runner and all of the commands have been started.
	StepRunning StepPhase = "Running"
	// StepSucceeded means that all containers in the Step have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	StepSucceeded StepPhase = "Succeeded"
	// StepFailed means that all containers in the Step have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	StepFailed StepPhase = "Failed"
	// StepUnknown means that for some reason the state of the Step could not be obtained, typically due
	// to an error in communicating with the host of the Step.
	StepUnknown StepPhase = "Unknown"
)

type Step struct {
	Id int32 `json:"id" protobuf:"varint,1,opt,name=id"`
	// The phase of a Step is a simple, high-level summary of where the Step is in its lifecycle.
	Phase StepPhase `json:"status" protobuf:"bytes,2,opt,name=status"`
	// Envs were the environment values which would be used by the called shell script.
	// Usually, they would include some base configuration
	Envs map[string]string `json:"envs" protobuf:"bytes,3,opt,name=envs"`
}
