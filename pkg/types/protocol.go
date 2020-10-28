package types

// +Protocol
// Request was the context which would be sent from the web dashboard or an abstract Runner
type Request struct {
	Type Type   `json:"type" protobuf:"bytes,1,opt,name=type"`
	Data []byte `json:"data" protobuf:"bytes,2,opt,name=data"`
}

// +Protocol
// Response was the context which would be sent from the Scheduler.
type Response struct {
	Type Type   `json:"type" protobuf:"bytes,1,opt,name=type"`
	Data []byte `json:"data" protobuf:"bytes,2,opt,name=data"`
}

type Body string

const (
	BodyRunner    Body = "Runner"
	BodyDashboard Body = "Dashboard"
)

// +Protocol
// Type
type Type struct {
	Namespace string  `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName string  `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	Body      Body    `json:"body" protobuf:"bytes,3,opt,name=body"`
	Command   Command `json:"command" protobuf:"bytes,4,opt,name=command"`
}

type Content string

const (
	ContentRunnerInfo Content = "RunnerInfo"
	ContentLogStream  Content = "LogStream"
)

// +Protocol
// RunnerInfo was the full information about a Runner, it would be sent from each remote abstract Runner.
type RunnerInfo struct {
	Name      string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Hostname  string `json:"hostname" protobuf:"bytes,2,opt,name=hostname"`
	Namespace string `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	GroupName string `json:"groupName" protobuf:"bytes,4,opt,name=groupName"`
	Steps     []Step `json:"steps" protobuf:"bytes,5,opt,name=steps"`
}

// +Protocol
// LogStream was the string which was transferred from the abstract Runner when the Runner was running a step.
// And it would also be sent from the Scheduler to each web dashboard for showing and watching
type LogStream struct {
	Namespace string `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName string `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	StepName  string `json:"stepName" protobuf:"bytes,3,opt,name=stepName"`
	Output    []byte `json:"output" protobuf:"bytes,4,opt,name=output"`
}

type Command string

const (
	CommandUpdate Command = "Update"
	CommandRun    Command = "Run"
	CommandResult Command = "Result"
)

// +Protocol
//
type StepCommand struct {
	Step Step `json:"step" protobuf:"bytes,1,opt,name=step"`
}
