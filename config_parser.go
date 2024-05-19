package main

import (
	"encoding/json"
	"os"
)

func ReadAndParseConfig(configPath string) (*Config, error) {
	// 打开配置文件
	configFile, err := os.Open(configPath)
	if err != nil {
		return &Config{}, err
	}
	defer configFile.Close()

	// 解析配置文件
	c := &Config{}
	err = json.NewDecoder(configFile).Decode(c)
	if err != nil {
		return c, err
	}
	return c, nil
}
