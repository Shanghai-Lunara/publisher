package main

import (
	"flag"
	"github.com/Shanghai-Lunara/publisher/pkg/conf"
	"github.com/Shanghai-Lunara/publisher/pkg/scheduler"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
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
	s := scheduler.NewServer(conf.Init())
	<-stopCh
	s.Close()
}
