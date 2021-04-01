package main

import (
	"flag"
	"github.com/Shanghai-Lunara/publisher/pkg/conf"
	"github.com/Shanghai-Lunara/publisher/pkg/scheduler"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
	"k8s.io/klog/v2"
)

func main() {
	var configPath = flag.String("configPath", "conf.yml", "configuration file path")
	klog.InitFlags(nil)
	flag.Parse()
	stopCh := signals.SetupSignalHandler()
	s := scheduler.NewServer(conf.Init(*configPath))
	<-stopCh
	s.Shutdown()
	<-stopCh
}
