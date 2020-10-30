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
	Namespace  string     `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName  string     `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	Body       Body       `json:"body" protobuf:"bytes,3,opt,name=body"`
	ServiceAPI ServiceAPI `json:"serviceApi" protobuf:"bytes,4,opt,name=serviceApi"`
}

type Content string

const (
	ContentRunnerInfo Content = "RunnerInfo"
	ContentLogStream  Content = "LogStream"
)

type RunnerType string

const (
	RunnerTypeClient RunnerType = "client"
	RunnerTypeServer RunnerType = "server"
)

// +Protocol
// RunnerInfo was the full information about a Runner, it would be sent from each remote abstract Runner.
type RunnerInfo struct {
	Name       string     `json:"name" protobuf:"bytes,1,opt,name=name"`
	Hostname   string     `json:"hostname" protobuf:"bytes,2,opt,name=hostname"`
	Namespace  string     `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	GroupName  string     `json:"groupName" protobuf:"bytes,4,opt,name=groupName"`
	RunnerType RunnerType `json:"runnerType" protobuf:"bytes,5,opt,name=runnerType"`
	Steps      []Step     `json:"steps" protobuf:"bytes,6,opt,name=steps"`
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

type ServiceAPI string

const (
	ListNamespace  ServiceAPI = "ListNamespace"
	ListGroupName  ServiceAPI = "ListGroupName"
	ListTask       ServiceAPI = "ListTask"
	RegisterRunner ServiceAPI = "RegisterRunner"
	UpdateStep     ServiceAPI = "UpdateStep"
	RunStep        ServiceAPI = "RunStep"
)

type Result struct {
	Items []string `json:"items" protobuf:"bytes,1,opt,name=items"`
}

type ListNamespaceRequest struct {
}

type ListNamespaceResponse struct {
	Items []string `json:"items" protobuf:"bytes,1,opt,name=items"`
}

type ListGroupNameRequest struct {
	Namespace Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
}

type ListGroupNameResponse struct {
	Items []string `json:"items" protobuf:"bytes,1,opt,name=items"`
}

type ListTaskRequest struct {
	Namespace Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName GroupName `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	Page      int32     `json:"page" protobuf:"varint,3,opt,name=page"`
	Length    int32     `json:"length" protobuf:"varint,4,opt,name=length"`
}

type ListTaskResponse struct {
	Tasks []Task `json:"tasks" protobuf:"bytes,1,opt,name=tasks"`
}
