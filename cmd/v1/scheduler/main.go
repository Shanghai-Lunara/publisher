package main

import (
	"flag"
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/Shanghai-Lunara/publisher/pkg/conf"
	"github.com/Shanghai-Lunara/publisher/pkg/scheduler"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
	"k8s.io/klog/v2"
)

func main() {
	fmt.Println(121212121)
	var configPath = flag.String("configPath", "conf.yml", "configuration file path")
	//var rbacPath = flag.String("rbacPath", "rbac_model.conf", "the default conf of the rbac mode")
	klog.InitFlags(nil)
	flag.Parse()
	defer zaplogger.Sync()
	stopCh := signals.SetupSignalHandler()
	s := scheduler.NewServer(conf.Init(*configPath), "")
	<-stopCh
	s.Shutdown()
	<-stopCh
}
