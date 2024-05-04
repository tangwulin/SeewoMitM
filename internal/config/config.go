package config

type Config struct {
	// 日志等级
	LogLevel string `json:"logLevel" default:"info"`

	// 屏幕保护劫持模式
	ScreenSaverHijackMode ScreenSaverHijackMode `json:"screenSaverHijackMode"`

	// 屏幕保护劫持内容
	ScreenSaverHijackContent []ScreenSaverHijackContent `json:"screenSaverHijackContent"`

	// 屏幕保护右下角来源显示
	ScreenSaverSource string `json:"screenSaverSource"`

	// 屏幕保护触发时间
	ScreenSaverEmitTime int `json:"screenSaverEmitTime" default:"600"`

	// MitM配置
	MitM *MitMConfig `json:"middleware,omitempty"`

	// 缓存配置
	Cache *CacheConfig `json:"cache,omitempty"`
}
