package helper

import "SeewoMitM/internal/config"

var globalConfig = config.Config{}

func GetConfig() *config.Config {
	return &globalConfig
}

func SetConfig(config config.Config) {
	globalConfig = config
}
