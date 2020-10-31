package scheduler

import (
	"context"
	"k8s.io/klog"
	"net"
	"net/http"
)

type Server struct {
	connections *connections
}

func (s *Server) initWSServer(addr string) {
	klog.Info("initWSService")
	http.HandleFunc("/runner", s.connections.handlerRunner)
	http.HandleFunc("/dashboard", s.connections.handlerDashboard)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Fatal(http.Serve(l, nil))
}

func (s *Server) Close() {

}

func NewServer(addr string) *Server {
	s := &Server{
		connections: NewConnections(context.Background()),
	}
	go s.initWSServer(addr)
	return s
}
