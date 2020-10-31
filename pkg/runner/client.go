package runner

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/nevercase/publisher/pkg/scheduler"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog"
	"net/url"
	"time"
)

type client struct {
	conn      *websocket.Conn
	writeChan chan []byte
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewClient(addr string) (*client, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: types.WebsocketHandlerRunner}
	klog.Info("url:", u)
	a, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &client{
		conn:      a,
		writeChan: make(chan []byte, 1024),
		ctx:       ctx,
		cancel:    cancel,
	}
	go c.readPump()
	go c.writePump()
	c.register()
	go c.ping()
	return c, nil
}

func (c *client) register() {
	req1 := &types.RegisterRunnerRequest{}
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

func (c *client) ping() {
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

func (c *client) readPump() {
	for {
		messageType, message, err := c.conn.ReadMessage()
		klog.Infof("messageType: %d message: %s err:%v\n", messageType, string(message), err)
		if err != nil {
			klog.Fatal(err)
			return
		}
	}
}

func (c *client) writePump() {
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
