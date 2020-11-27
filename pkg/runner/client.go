package runner

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/nevercase/publisher/pkg/scheduler"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog/v2"
	"net/url"
	"time"
)

type Client struct {
	conn         *websocket.Conn
	writeChan    chan []byte
	runner       *Runner
	streamOutput chan string
	pingTimer    bool
	currentStep  *types.Step
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewClient(addr string, streamOutput chan string, r *Runner) (*Client, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: types.WebsocketHandlerRunner}
	klog.Info("url:", u)
	a, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		klog.Fatal(err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		conn:         a,
		writeChan:    make(chan []byte, 1024),
		runner:       r,
		streamOutput: streamOutput,
		pingTimer:    false,
		ctx:          ctx,
		cancel:       cancel,
	}
	c.runner.StreamOutput = c.streamOutput
	go c.readPump()
	go c.writePump()
	go c.logStream()
	c.register()
	return c, nil
}

func (c *Client) register() {
	ri, err := c.runner.Register()
	if err != nil {
		klog.Fatal(err)
	}
	req1 := &types.RegisterRunnerRequest{
		RunnerInfo: ri,
	}
	data, err := req1.Marshal()
	if err != nil {
		klog.Fatal(err)
	}
	req2 := &types.Request{
		Type: types.Type{
			Body:       types.BodyRunner,
			ServiceAPI: types.RegisterRunner,
		},
		Data: data,
	}
	data2, err := req2.Marshal()
	if err != nil {
		klog.Fatal(err)
	}
	c.writeChan <- data2
}

func (c *Client) ping() {
	tick := time.NewTicker(time.Second * time.Duration(scheduler.WebsocketConnectionTimeout/2))
	defer tick.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-tick.C:
			var data []byte
			req := &types.Request{
				Type: types.Type{
					Body:       types.BodyRunner,
					ServiceAPI: types.Ping,
				},
				Data: data,
			}
			res, err := req.Marshal()
			if err != nil {
				klog.Fatal(err)
				return
			}
			c.writeChan <- res
		}
	}
}

func (c *Client) readPump() {
	for {
		messageType, message, err := c.conn.ReadMessage()
		klog.V(5).Infof("messageType: %d message: %s err:%v\n", messageType, string(message), err)
		if err != nil {
			klog.Fatal(err)
			return
		}
		req := &types.Request{}
		if err := req.Unmarshal(message); err != nil {
			klog.Fatal(err)
			return
		}

		switch req.Type.ServiceAPI {
		case types.RegisterRunner:
			if c.pingTimer == false {
				go c.ping()
			}
		case types.Ping:
		case types.RunStep:
			data := &types.RunStepRequest{}
			if err = data.Unmarshal(req.Data); err != nil {
				klog.Fatal(err)
			}
			go func() {
				c.currentStep = &data.Step
				if err = c.runner.Run(&data.Step); err != nil {
					klog.V(2).Info(err)
					// todo catching error, update Step's Messages, and report to Scheduler
				}
				if err = c.updateStepInformationToScheduler(&data.Step); err != nil {
					klog.Fatal(err)
				}
			}()

		case types.UpdateStep:
			data := &types.UpdateStepRequest{}
			if err = data.Unmarshal(req.Data); err != nil {
				klog.Fatal(err)
			}
			if err = c.runner.Update(&data.Step); err != nil {
				klog.V(2).Info(err)
				// todo catching error, update Step's Messages, and report to Scheduler
			}
		}
	}
}

func (c *Client) writePump() {
	for {
		select {
		case msg, isClose := <-c.writeChan:
			if !isClose {
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				klog.Fatal(err)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) logStream() {
	for {
		select {
		case log, isClose := <-c.streamOutput:
			if !isClose {
				return
			}
			if c.currentStep == nil {
				klog.V(2).Info("Client currentStep was nil")
				continue
			}
			req1 := &types.LogStreamRequest{
				Namespace:  c.runner.Namespace,
				GroupName:  c.runner.GroupName,
				RunnerName: c.runner.Name,
				StepName:   c.currentStep.Name,
				Output:     log,
			}
			data, err := req1.Marshal()
			if err != nil {
				klog.Fatal(err)
			}
			req2 := &types.Request{
				Type: types.Type{
					Body:       types.BodyRunner,
					ServiceAPI: types.LogStream,
				},
				Data: data,
			}
			data, err = req2.Marshal()
			if err != nil {
				klog.Fatal(err)
			}
			c.writeChan <- data
		}
	}
}

func (c *Client) updateStepInformationToScheduler(s *types.Step) (err error) {
	s, err = c.runner.Step(s)
	if err != nil {
		klog.Fatal(err)
	}
	req1 := &types.UpdateStepRequest{
		Namespace:  c.runner.Namespace,
		GroupName:  c.runner.GroupName,
		RunnerName: c.runner.Name,
		Step:       *s,
	}
	data, err := req1.Marshal()
	if err != nil {
		klog.Fatal(err)
	}
	req2 := &types.Request{
		Type: types.Type{
			Body:       types.BodyRunner,
			ServiceAPI: types.UpdateStep,
		},
		Data: data,
	}
	data, err = req2.Marshal()
	if err != nil {
		klog.Fatal(err)
	}
	c.writeChan <- data
	return nil
}
