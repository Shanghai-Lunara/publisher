package conf

import (
	"flag"
	"github.com/nevercase/publisher/pkg/dao"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog"
)

type GameService struct {
	ListenPort int `json:"listenPort" yaml:"listenPort"`
}

type Config struct {
	GameService GameService         `yaml:"GameService,flow"`
	Mysql       dao.MysqlPoolConfig `yaml:"Mysql,flow"`
}

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "conf.yml", "configuration file path")
}

func Init() (*Config, error) {
	c := &Config{}
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(configPath); err != nil {
		klog.Fatal(err)
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		klog.Fatal(err)
	}
	return c, nil
}