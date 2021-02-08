package conf

import (
	"github.com/Shanghai-Lunara/publisher/pkg/dao"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog"
)

type PublisherService struct {
	ListenPort int `json:"listenPort" yaml:"listenPort"`
}

type Config struct {
	PublisherService PublisherService    `yaml:"PublisherService,flow"`
	Mysql            dao.MysqlPoolConfig `yaml:"Mysql,flow"`
	Projects         []Project           `yaml:"Projects"`
}

type Project struct {
	Namespace string  `yaml:"namespace"`
	Groups    []Group `yaml:"groups"`
}

type Group struct {
	Name string `yaml:"name"`
}

func Init(file string) *Config {
	c := &Config{}
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(file); err != nil {
		klog.Fatal(err)
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		klog.Fatal(err)
	}
	return c
}
