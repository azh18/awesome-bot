package lark

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type BotConfig struct {
	AppID       string `yaml:"AppId"`
	AppSecret   string `yaml:"AppSecret"`
	VerifyToken string `yaml:"VerifyToken"`
	EncryptKey  string `yaml:"EncryptKey"`
}

type User struct {
	Name   string `yaml:"Name"`
	ChatID string `yaml:"ChatID"`
	Email  string `yaml:"Email"`
}

type Config struct {
	BotConfig *BotConfig `yaml:"BotConfig"`
	Users     []*User    `yaml:"Users"` // name -> chatID
}

func ParseConfig(confPath string) *Config {
	configFile, err := os.Open(confPath)
	if err != nil {
		klog.Fatalf("open file %s error: %s", confPath, err.Error())
	}
	content, err := ioutil.ReadAll(configFile)
	if err != nil {
		klog.Fatalf("read file %s error: %s", confPath, err.Error())
	}
	config := &Config{}
	if err := yaml.Unmarshal(content, config); err != nil {
		klog.Fatalf("unmarshal lark config file error: %s", err.Error())
	}
	return config
}
