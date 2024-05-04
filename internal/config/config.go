package config

type Config struct {
	// 日志等级
	LogLevel string `json:"logLevel" default:"info"`

	// 屏幕保护相关配置
	ScreensaverConfig *ScreensaverConfig `json:"screensaverConfig"`
}
