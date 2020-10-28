package scheduler

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/klog"
)

const (
	WebsocketConnectionTimeout = 10
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewConnections(ctx context.Context) *connections {
	return &connections{
		autoIncrementId: 0,
		items:           make(map[int32]*conn, 0),
		ctx:             ctx,
	}
}

type connections struct {
	mu              sync.Mutex
	autoIncrementId int32
	items           map[int32]*conn
	ctx             context.Context
}

func (cs *connections) handler(w http.ResponseWriter, r *http.Request) {
	c, err := cs.newConn(w, r)
	if err != nil {
		klog.V(2).Info(err)
		return
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.items[c.id] = c
}

func (cs *connections) newConn(w http.ResponseWriter, r *http.Request) (*conn, error) {
	client, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(cs.ctx)
	c := &conn{
		id:                    atomic.AddInt32(&cs.autoIncrementId, 1),
		conn:                  client,
		writeChan:             make(chan []byte, 4096),
		lastPingTime:          time.Now(),
		keepAliveTimeoutInSec: WebsocketConnectionTimeout,
		closeOnce:             sync.Once{},
		ctx:                   ctx,
		cancel:                cancel,
	}
	go c.keepAlive()
	go c.readPump()
	go c.writePump()
	return c, nil
}

type conn struct {
	id                    int32
	conn                  *websocket.Conn
	writeChan             chan []byte
	lastPingTime          time.Time
	keepAliveTimeoutInSec int64
	closeOnce             sync.Once
	ctx                   context.Context
	cancel                context.CancelFunc
}

func (c *conn) keepAlive() {
	defer c.close()
	tick := time.NewTicker(time.Second * time.Duration(c.keepAliveTimeoutInSec+1))
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if time.Now().Sub(c.lastPingTime) > time.Second*time.Duration(c.keepAliveTimeoutInSec) {
				klog.Info("keepAlive timeout")
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *conn) close() {
	c.closeOnce.Do(func() {
		c.cancel()
	})
}

func (c *conn) readPump() {
	defer c.close()
	for {
		messageType, data, err := c.conn.ReadMessage()
		klog.Infof("messageType: %d message-string: %v\n", messageType, string(data))
		if err != nil {
			klog.V(2).Info(err)
			return
		}
		// todo handle message
	}
}

func (c *conn) writePump() {
	defer c.close()
	for {
		select {
		case <-c.ctx.Done():
			return
		}
	}
}
