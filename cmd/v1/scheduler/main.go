package main

import (
	"flag"

	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
	"github.com/nevercase/publisher/pkg/scheduler"
	"k8s.io/klog/v2"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":6969", "The address of the Publisher.")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	stopCh := signals.SetupSignalHandler()
	s := scheduler.NewServer(addr)
	<-stopCh
	s.Close()
}
