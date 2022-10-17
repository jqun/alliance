package config

import (
	"alliance/common/consts"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	App *AppConfig `toml:"App"`
}

var serverConfig *ServerConfig

func init() {
	serverConfig = &ServerConfig{}
	serverConfig.loadData()
}

func (m *ServerConfig) loadData() {
	path := getConfigPath()
	_, err := toml.DecodeFile(path, m)
	if err != nil {
		log.Panicf("load config data err[%v]", err.Error())
	}
}

func getConfigPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, consts.Dir, "server.toml")
}

func GetServerConfig() *ServerConfig {
	return serverConfig
}
