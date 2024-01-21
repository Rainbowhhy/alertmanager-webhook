package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	QyWechatKey   string `yaml:"qywechatKey"`
	FeishuKey     string `yaml:"feishuKey"`
	DingdingKey   string `yaml:"dingdingKey"`
	LogFileDir    string `yaml:"logFileDir"`
	LogFilePath   string `yaml:"logFilePath"`
	Port          string `yaml:"port"`
	Host          string `yaml:"host"`
	RedisServer   string `yaml:"redisServer"`
	RedisPort     string `yaml:"redisPort"`
	RedisPassword string `yaml:"redisPassword"`
}

func GetConf(workdir, configFile string) Config {
	var webhookConf Config
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("The program configuration file does not exist, please configure it correctly!")
		os.Exit(100)
	}

	// 读取配置文件，从yaml解析为struct
	err = yaml.Unmarshal(yamlFile, &webhookConf)
	if err != nil {
		fmt.Println(err.Error())
	}
	if webhookConf.LogFileDir != "" {
		webhookConf.LogFilePath = path.Join(webhookConf.LogFileDir, webhookConf.LogFilePath)
	} else {
		webhookConf.LogFilePath = path.Join(workdir, webhookConf.LogFilePath)
	}
	if webhookConf.Port == "" {
		webhookConf.Port = "9095"
	}
	if webhookConf.Host == "" {
		webhookConf.Host = "127.0.0.1"
	}
	if webhookConf.RedisPort == "" {
		webhookConf.RedisPort = "6379"
	}

	return webhookConf
}
