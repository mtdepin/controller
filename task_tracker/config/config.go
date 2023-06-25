package config

import (
	"controller/api"
	"controller/pkg/config"
	"controller/pkg/db"
	"controller/pkg/logger"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
)

type ServerConfig struct {
	Node      NodeConfig
	Logger    config.LoggerConfig
	Jeager    config.JaegerConfig
	Request   config.RequestConfig
	Strategy  config.StrategyConfig
	Scheduler config.SchedulerConfig
	Montior   config.MontiorConfig
	DB        db.DBconfig
}

var ServerCfg *ServerConfig

type NodeConfig struct {
	config.CommonNodeConfig
	NameServer_group string
	Api              string
}

func LoadServerConfig(serviceId string) (*ServerConfig, error) {
	if err := initConfig(serviceId, ""); err != nil {
		return nil, err
	}
	var c ServerConfig
	err := config.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	fmt.Println(c)

	//todo： 由于上面不能解析出一些字段，这里单独解析通用node配置，
	var cc config.CommonNodeConfig
	err = config.UnmarshalKey("node", &cc)
	if err != nil {
		return nil, err
	}
	c.Node.CommonNodeConfig = cc

	if err := c.checkConfig(); err != nil {
		return nil, err
	}
	return &c, nil
}

func initConfig(serviceId string, cmdRoot string) error {
	configInstance := config.InitViper(serviceId)
	defer func() {
		configInstance.WatchConfig()
		configInstance.Config.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("配置发生变更：", e.Name)
		})
	}()
	curPath, _ := os.Getwd()
	confPath := curPath + "/conf/"
	configInstance.AddConfigPath(confPath)
	if !api.FileExists(confPath + "task_tracker.yml") {
		return errors.New(confPath + "task_tracker.yml do not exist")
	}
	configInstance.SetConfigName(serviceId)
	if err := configInstance.ReadInConfig(); err != nil {
		err := fmt.Errorf("ctlerror when reading %s.json config file %s", cmdRoot, err)
		logger.Error(err)
		return err
	}
	return nil
}

func (c *ServerConfig) checkConfig() error {
	if c.Node.Node_group == "" {
		c.Node.Node_group = "ck_group"
	}

	if c.Node.NameServer_group == "" {
		c.Node.NameServer_group = "ns_group"
	}

	if c.Node.Region == "" {
		c.Node.Region = "cd"
	}

	if c.Node.Api == "" {
		c.Node.Api = "127.0.0.1:8521"
	}

	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}

	if c.Request.Max == 0 {
		c.Request.Max = 1024
	}

	if c.Request.TimeOut == 0 {
		c.Request.Max = 120
	}

	return nil

}
