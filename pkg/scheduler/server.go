package scheduler

import (
	"context"
	"fmt"
	"github.com/Shanghai-Lunara/publisher/pkg/conf"
	"github.com/Shanghai-Lunara/publisher/pkg/types"
	"k8s.io/klog/v2"
	"net"
	"net/http"
)

type Server struct {
	connections *connections
}

func (s *Server) initWSServer(addr string) {
	klog.Info("initWSService")
	http.HandleFunc(types.WebsocketHandlerRunner, s.connections.handlerRunner)
	http.HandleFunc(types.WebsocketHandlerDashboard, s.connections.handlerDashboard)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Fatal(http.Serve(l, nil))
}

func (s *Server) Close() {

}

func NewServer(c *conf.Config) *Server {
	s := &Server{
		connections: NewConnections(context.Background(), c),
	}
	go s.initWSServer(fmt.Sprintf(":%d", c.PublisherService.ListenPort))
	return s
}
