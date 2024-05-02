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
}

type ScreenSaverHijackMode int32

type ScreenSaverHijackContent struct {
	Type           ScreenSaverHijackContentType `json:"type"`
	Path           string                       `json:"path,omitempty"`
	RequirePreload bool                         `json:"requirePreload,omitempty"`
	EntryPoint     string                       `json:"entryPoint,omitempty"`
	ServePath      string                       `json:"servePath,omitempty"`
	SpineVersion   string                       `json:"spineVersion,omitempty"`
	SpineConfig    interface{}                  `json:"spineConfig,omitempty"`
}

const (
	ScreenSaverHijackModeOff ScreenSaverHijackMode = iota
	ScreenSaverHijackModeAdd
	ScreenSaverHijackModeReplaceAll
)

type ScreenSaverHijackContentType string

const (
	HTMLScreenSaverHijackContent           ScreenSaverHijackContentType = "html"
	ImageScreenSaverHijackContent          ScreenSaverHijackContentType = "image"
	VideoScreenSaverHijackContent          ScreenSaverHijackContentType = "video"
	SpineScreenSaverHijackContent          ScreenSaverHijackContentType = "spine"
	ImageDirectoryScreenSaverHijackContent ScreenSaverHijackContentType = "imageDirectory"
)
