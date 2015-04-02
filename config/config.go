package config

import (
	"fmt"
	"io/ioutil"

	goyaml "gopkg.in/yaml.v2"
)

type Config struct {
	MasterOp struct {
		Rdb bool
		Aof bool
	}
	ConnErrLlimit   int
	ConcurrentLimit int
	ProxyId         string
	ProductName     string
	ZkAddr          string
}

var ProxyConfig Config

func ReloadConfig(path string) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = goyaml.Unmarshal(content, &ProxyConfig); err != nil {
		return
	}
	fmt.Println(ProxyConfig)
	return
}

func init() {
	//config_path := flag.String("config", "/etc/xiaoju/codisconfig.yaml", "config file's path")
	//flag.Parse()
	//if err := ReloadConfig(*config_path); err != nil {
	//	panic(fmt.Sprintf("reload config failed, because of %v", err))
	//}
}
