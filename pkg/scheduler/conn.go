package scheduler

import (
	"context"
	"github.com/nevercase/publisher/pkg/conf"
	"github.com/nevercase/publisher/pkg/dao"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog/v2"
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

func NewConnections(ctx context.Context, c *conf.Config) *connections {
	cs := &connections{
		autoIncrementId: 0,
		items:           make(map[int32]*conn, 0),
		broadcast:       make(chan *broadcast, 1024),
		removedChan:     make(chan int32, 100),
		ctx:             ctx,
	}
	cs.scheduler = NewScheduler(cs.broadcast, dao.New(&c.Mysql))
	go cs.remove()
	go cs.broadcastToDashboard()
	return cs
}

type connections struct {
	mu              sync.RWMutex
	autoIncrementId int32
	items           map[int32]*conn
	broadcast       chan *broadcast
	removedChan     chan int32
	scheduler       *Scheduler
	ctx             context.Context
}

type broadcastType string

const (
	broadcastTypeBindRunner broadcastType = "bindRunner"
	broadcastTypePing       broadcastType = "ping"
	broadcastTypeDashboard  broadcastType = "dashboard"
	broadcastTypeRunner     broadcastType = "runner"
)

type broadcast struct {
	bt         broadcastType
	clientId   int32
	runnerName string
	msg        []byte
}

func (cs *connections) broadcastToDashboard() {
	for {
		select {
		case broadcast, isClose := <-cs.broadcast:
			if !isClose {
				return
			}
			switch broadcast.bt {
			case broadcastTypeBindRunner:
				cs.mu.RLock()
				if t, ok := cs.items[broadcast.clientId]; ok {
					t.runnerName = broadcast.runnerName
				} else {
					// todo handle err
				}
				cs.mu.RUnlock()
			case broadcastTypePing:
				cs.mu.RLock()
				if t, ok := cs.items[broadcast.clientId]; ok {
					t.ping()
				} else {
					// todo handle err
				}
				cs.mu.RUnlock()
			case broadcastTypeDashboard:
				cs.mu.RLock()
				for _, v := range cs.items {
					if v.body == types.BodyDashboard {
						v.writeChan <- broadcast.msg
					}
				}
				cs.mu.RUnlock()
			case broadcastTypeRunner:
				cs.mu.RLock()
				for _, v := range cs.items {
					if v.runnerName == broadcast.runnerName {
						v.writeChan <- broadcast.msg
					}
				}
				cs.mu.RUnlock()
			}
		case <-cs.ctx.Done():
			return
		}
	}
}

func (cs *connections) remove() {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case id, isClose := <-cs.removedChan:
			if !isClose {
				return
			}
			cs.mu.Lock()
			delete(cs.items, id)
			cs.mu.Unlock()
		}
	}
}

func (cs *connections) handlerDashboard(w http.ResponseWriter, r *http.Request) {
	c, err := cs.newConn(w, r, types.BodyDashboard)
	if err != nil {
		klog.V(2).Info(err)
		return
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.items[c.id] = c
}

func (cs *connections) handlerRunner(w http.ResponseWriter, r *http.Request) {
	c, err := cs.newConn(w, r, types.BodyRunner)
	if err != nil {
		klog.V(2).Info(err)
		return
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.items[c.id] = c
}

func (cs *connections) newConn(w http.ResponseWriter, r *http.Request, body types.Body) (*conn, error) {
	client, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(cs.ctx)
	c := &conn{
		scheduler:             cs.scheduler,
		body:                  body,
		id:                    atomic.AddInt32(&cs.autoIncrementId, 1),
		conn:                  client,
		writeChan:             make(chan []byte, 4096),
		lastPingTime:          time.Now(),
		keepAliveTimeoutInSec: WebsocketConnectionTimeout,
		closeOnce:             sync.Once{},
		removedChan:           cs.removedChan,
		ctx:                   ctx,
		cancel:                cancel,
	}
	go c.keepAlive()
	go c.readPump()
	go c.writePump()
	return c, nil
}

// conn was an abstract runner or a web dashboard client
type conn struct {
	scheduler             *Scheduler
	body                  types.Body
	id                    int32
	runnerName            string
	conn                  *websocket.Conn
	writeChan             chan []byte
	lastPingTime          time.Time
	keepAliveTimeoutInSec int64
	closeOnce             sync.Once
	removedChan           chan<- int32
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

func (c *conn) ping() {
	c.lastPingTime = time.Now()
}

func (c *conn) close() {
	c.closeOnce.Do(func() {
		c.cancel()
		c.removedChan <- c.id
		if err := c.conn.Close(); err != nil {
			klog.V(2).Info(err)
		}
	})
}

func (c *conn) readPump() {
	defer c.close()
	for {
		messageType, data, err := c.conn.ReadMessage()
		klog.V(5).Info("data:", data)
		klog.V(5).Infof("messageType: %d message-string: %s\n", messageType, string(data))
		if err != nil {
			klog.V(2).Info(err)
			return
		}
		res, err := c.scheduler.handle(data, c.id)
		if err != nil {
			return
		}
		if len(res) > 0 {
			c.writeChan <- res
		}
	}
}

func (c *conn) writePump() {
	defer c.close()
	for {
		select {
		case msg, isClose := <-c.writeChan:
			if !isClose {
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				klog.V(2).Info(err)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}
