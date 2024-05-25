package main

import (
	"SeewoMitM/model"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
)

type Config struct {
	// 日志等级
	LogLevel string `json:"logLevel" default:"info"`

	// 屏幕保护相关配置
	ScreensaverConfig *model.ScreensaverConfig `json:"screensaverConfig"`
}

var globalConfig *Config

func InitConfig(configPath string) error {
	dir, fileName := path.Split(configPath)
	viper.SetConfigFile(fileName)
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("config file not found, creating default config file...")
			// 创建默认配置文件
			file, err := os.Create(configPath)
			if err != nil {
				fmt.Printf("create config file error: %v\n", err)
				return err
			}
			defer file.Close()

			viper.SetDefault("logLevel", "info")
			viper.SetDefault("screensaverConfig", *NewScreensaverConfig())
		} else {
			fmt.Println("read config file error:", err)
			return err
		}

		err := viper.WriteConfig()
		if err != nil {
			return err
		}
	}

	if err := viper.Unmarshal(&globalConfig); err != nil {
		fmt.Println("unmarshal config file error:", err)
		return err
	}

	return nil
}

func SaveConfig() error {
	return viper.WriteConfig()
}
