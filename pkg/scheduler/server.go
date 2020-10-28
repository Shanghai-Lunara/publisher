package scheduler

import (
	"k8s.io/klog"
	"net"
	"net/http"
)

type Server struct {
	connections *connections
}

func (s *Server) initWSServer(addr string) {
	klog.Info("initWSService")
	http.HandleFunc("/", s.connections.handler)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Fatal(http.Serve(l, nil))
}
