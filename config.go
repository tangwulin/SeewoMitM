package main

import "SeewoMitM/model"

type Config struct {
	// 日志等级
	LogLevel string `json:"logLevel" default:"info"`

	// 屏幕保护相关配置
	ScreensaverConfig *model.ScreensaverConfig `json:"screensaverConfig"`
}

var globalConfig = Config{}

func GetConfig() *Config {
	return &globalConfig
}

func SetConfig(config Config) {
	globalConfig = config
}
