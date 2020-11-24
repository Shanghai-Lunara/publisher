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

type StepPolicy string

const (
	StepPolicyAuto   StepPolicy = "auto"
	StepPolicyManual StepPolicy = "manual"
)

type StepAvailable string

const (
	StepAvailableEnable  StepAvailable = "enable"
	StepAvailableDisable StepAvailable = "disable"
)

type Step struct {
	Id int32 `json:"id" protobuf:"varint,1,opt,name=id"`
	// Name was the name of a Step which must be unique
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
	// The phase of a Step is a simple, high-level summary of where the Step is in its lifecycle.
	Phase StepPhase `json:"status" protobuf:"bytes,3,opt,name=status"`
	// Policy was the StepPolicy of the Step which could control the Runner
	Policy StepPolicy `json:"policy" protobuf:"bytes,4,opt,name=policy"`
	// Available determines whether the Step was available, if the value was disable that means the current
	// step should be skipped instead of running.
	Available StepAvailable `json:"available" protobuf:"bytes,5,opt,name=available"`
	// Envs were the environment values which would be used by the called shell script.
	// Usually, they would include some base configuration
	Envs map[string]string `json:"envs" protobuf:"bytes,6,opt,name=envs"`
	// Output was the stdout from the executing shell commands
	Output []string `json:"output" protobuf:"bytes,7,opt,name=output"`
	// UploadFiles were the map of the files which would be uploaded to the remote ftp server by the Step Run().
	UploadFiles []UploadFile `json:"uploadFiles" protobuf:"bytes,8,opt,name=uploadFiles"`
	// WriteFiles were map of the files which would be written to the remote ftp server by the Step Run().
	WriteFiles []WriteFile `json:"writeFiles" protobuf:"bytes,9,opt,name=writeFiles"`
	// Messages were the human readable message indicating details about all the conditions in a Step's lifecycle.
	Messages []string `json:"messages" protobuf:"bytes,10,opt,name=messages"`
	// RunnerName was the name of a Runner (the Runner which has been called to run this Step)
	RunnerName string `json:"runnerName" protobuf:"bytes,11,opt,name=runnerName"`
	// DurationInMS was the total used time during the whole step
	DurationInMS int32 `json:"durationInMs" protobuf:"varint,12,opt,name=durationInMs"`
}

type UploadFile struct {
	// SourceFile was the absolute path about the file which has been marked to be uploaded later.
	SourceFile string `json:"sourceFile" protobuf:"bytes,1,opt,name=sourceFile"`
	// TargetPath was the target directory, it may be needed to be created before uploading the SourceFile
	TargetPath string `json:"targetPath" protobuf:"bytes,2,opt,name=targetPath"`
	// TargetFile was the relative path about the file which would be created in the ftp server.
	TargetFile string `json:"targetFile" protobuf:"bytes,3,opt,name=targetFile"`
}

type WriteFile struct {
	// Content was the context of the file which needed to be written to the remote ftp server.
	Content []byte `json:"content" protobuf:"bytes,1,opt,name=content"`
	// TargetFile was the relative path about the file which would be created in the ftp server.
	TargetFile string `json:"targetFile" protobuf:"bytes,2,opt,name=targetFile"`
}
