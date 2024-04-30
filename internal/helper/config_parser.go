package helper

import (
	"SeewoMitM/internal/config"
	"encoding/json"
	"os"
)

func ReadAndParseConfig(configPath string) (*config.Config, error) {
	// 打开配置文件
	configFile, err := os.Open(configPath)
	if err != nil {
		return &config.Config{}, err
	}
	defer configFile.Close()

	// 解析配置文件
	c := &config.Config{}
	err = json.NewDecoder(configFile).Decode(c)
	if err != nil {
		return c, err
	}
	return c, nil
}
