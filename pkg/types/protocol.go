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
	Code    int32  `json:"code" protobuf:"varint,1,opt,name=code"`
	Message string `json:"message" protobuf:"bytes,2,opt,name=message"`
	Type    Type   `json:"type" protobuf:"bytes,3,opt,name=type"`
	Data    []byte `json:"data" protobuf:"bytes,4,opt,name=data"`
}

type Body string

const (
	BodyRunner    Body = "Runner"
	BodyDashboard Body = "Dashboard"
)

// +Protocol
// Type
type Type struct {
	Body       Body       `json:"body" protobuf:"bytes,1,opt,name=body"`
	ServiceAPI ServiceAPI `json:"serviceApi" protobuf:"bytes,2,opt,name=serviceApi"`
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

type Namespace string

type GroupName string

type Group struct {
	Tasks   []Task       `json:"tasks" protobuf:"bytes,1,rep,name=tasks"`
	Runners []RunnerInfo `json:"runners" protobuf:"bytes,2,rep,name=runners"`
}

type Task struct {
	Id      int32                 `json:"id" protobuf:"varint,1,opt,name=id"`
	Runners map[string]RunnerInfo `json:"runners" protobuf:"bytes,2,opt,name=runners"`
}

// +Protocol
// RunnerInfo was the full information about a Runner, it would be sent from each remote abstract Runner.
type RunnerInfo struct {
	Name       string     `json:"name" protobuf:"bytes,1,opt,name=name"`
	Hostname   string     `json:"hostname" protobuf:"bytes,2,opt,name=hostname"`
	Namespace  Namespace  `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	GroupName  GroupName  `json:"groupName" protobuf:"bytes,4,opt,name=groupName"`
	RunnerType RunnerType `json:"runnerType" protobuf:"bytes,5,opt,name=runnerType"`
	Steps      []Step     `json:"steps" protobuf:"bytes,6,opt,name=steps"`
}

type ServiceAPI string

const (
	Ping                          ServiceAPI = "Ping"
	ListNamespace                 ServiceAPI = "ListNamespace"
	ListGroupName                 ServiceAPI = "ListGroupName"
	ListRunner                    ServiceAPI = "ListRunner"
	RegisterRunner                ServiceAPI = "RegisterRunner"
	UpdateStep                    ServiceAPI = "UpdateStep"
	RunStep                       ServiceAPI = "RunStep"
	LogStream                     ServiceAPI = "LogStream"
	CompleteStep                  ServiceAPI = "CompleteStep"
	ServiceAPIListRecordsRequest  ServiceAPI = "ListRecordsRequest"
	ServiceAPIListRecordsResponse ServiceAPI = "ListRecordsResponse"
)

type Result struct {
	Items []string `json:"items" protobuf:"bytes,1,opt,name=items"`
}

type PingRequest struct {
}

type PongResponse struct {
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

type ListRunnerRequest struct {
	Namespace Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName GroupName `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
}

type ListRunnerResponse struct {
	Runners []RunnerInfo `json:"runners" protobuf:"bytes,1,opt,name=runners"`
}

type RegisterRunnerRequest struct {
	RunnerInfo RunnerInfo `json:"runnerInfo" protobuf:"bytes,1,opt,name=runnerInfo"`
}

type RegisterRunnerResponse struct {
}

type RunStepRequest struct {
	Namespace  Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName  GroupName `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	RunnerName string    `json:"runnerName" protobuf:"bytes,3,opt,name=runnerName"`
	Step       Step      `json:"step" protobuf:"bytes,4,opt,name=step"`
}

type RunStepResponse struct {
}

type UpdateStepRequest struct {
	Namespace  Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName  GroupName `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	RunnerName string    `json:"runnerName" protobuf:"bytes,3,opt,name=runnerName"`
	Step       Step      `json:"step" protobuf:"bytes,4,opt,name=step"`
}

type UpdateStepResponse struct {
}

type CompleteStepRequest struct {
	Namespace  Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName  GroupName `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	RunnerName string    `json:"runnerName" protobuf:"bytes,3,opt,name=runnerName"`
	Step       Step      `json:"step" protobuf:"bytes,4,opt,name=step"`
}

type CompleteStepResponse struct {
}

// +Protocol
// LogStreamRequest was the string which was transferred from the abstract Runner when the Runner was running a step.
// And it would also be sent from the Scheduler to each web dashboard for showing and watching
type LogStreamRequest struct {
	Namespace  Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName  GroupName `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	RunnerName string    `json:"runnerName" protobuf:"bytes,3,opt,name=runnerName"`
	StepName   string    `json:"stepName" protobuf:"bytes,4,opt,name=stepName"`
	Output     string    `json:"output" protobuf:"bytes,5,opt,name=output"`
}

type LogStreamResponse struct {
}

// ListRecordsRequest
type ListRecordsRequest struct {
	Namespace  Namespace `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	GroupName  GroupName `json:"groupName" protobuf:"bytes,2,opt,name=groupName"`
	RunnerName string    `json:"runnerName" protobuf:"bytes,3,opt,name=runnerName"`
	Page       int32     `json:"page" protobuf:"varint,4,opt,name=page"`
	Length     int32     `json:"length" protobuf:"varint,5,opt,name=length"`
}

// ListRecordsResponse
type ListRecordsResponse struct {
	Params  ListRecordsRequest `json:"params" protobuf:"bytes,1,opt,name=namespace"`
	Records []Record           `json:"records" protobuf:"bytes,2,opt,name=records"`
}
