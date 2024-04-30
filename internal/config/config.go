package config

type Config struct {
	// 日志等级
	LogLevel string `json:"logLevel" default:"info"`

	// 屏幕保护劫持模式
	ScreenSaverHijackMode ScreenSaverHijackMode `json:"screenSaverHijackMode"`

	// 屏幕保护劫持内容
	ScreenSaverHijackContent []ScreenSaverHijackContent `json:"screenSaverHijackContent"`

	// 屏幕保护触发时间
	ScreenSaverEmitTime int `json:"screenSaverEmitTime" default:"600"`
}

type ScreenSaverHijackMode int32

const (
	Off ScreenSaverHijackMode = iota
	Add
	ReplaceAll
)

type ScreenSaverHijackContent struct {
	ScreenSaverHijackHTMLContent
	ScreenSaverHijackImageContent
	ScreenSaverHijackVideoContent
	ScreenSaverHijackSpineContent
	ScreenSaverHijackImageDirectoryContent
}

type ScreenSaverHijackContentType string

const (
	HTML           ScreenSaverHijackContentType = "html"
	Image          ScreenSaverHijackContentType = "image"
	Video          ScreenSaverHijackContentType = "video"
	Spine          ScreenSaverHijackContentType = "spine"
	ImageDirectory ScreenSaverHijackContentType = "imageDirectory"
)

type ScreenSaverHijackHTMLContent struct {
	Type       ScreenSaverHijackContentType `json:"type" default:"html"`
	EntryPoint string                       `json:"entryPoint"`
	ServePath  string                       `json:"servePath"`
}

type ScreenSaverHijackImageContent struct {
	Type           ScreenSaverHijackContentType `json:"type" default:"image"`
	Path           string                       `json:"path"`
	RequirePreload bool                         `json:"requirePreload" default:"true"`
}

type ScreenSaverHijackVideoContent struct {
	Type           ScreenSaverHijackContentType `json:"type" default:"video"`
	Path           string                       `json:"path"`
	RequirePreload bool                         `json:"requirePreload" default:"true"`
}

type ScreenSaverHijackSpineContent struct {
	Type         ScreenSaverHijackContentType `json:"type" default:"spine"`
	SpineVersion string                       `json:"spineVersion" default:"3.8"`
	SpineConfig  string                       `json:"spineConfig"`
}

type ScreenSaverHijackImageDirectoryContent struct {
	Type ScreenSaverHijackContentType `json:"type" default:"imageDirectory"`
	Path string                       `json:"path"`
}
