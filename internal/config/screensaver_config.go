package config

type ScreensaverConfig struct {
	// 屏幕保护劫持模式
	HijackMode string `json:"hijackMode"`

	// 屏幕保护劫持内容
	Contents []ScreensaverContent `json:"contents"`

	// 屏幕保护右下角来源显示
	Source string `json:"source"`

	// 屏幕保护触发时间
	EmitTime int `json:"emitTime" default:"600"`
}

func NewScreensaverConfig() *ScreensaverConfig {
	return &ScreensaverConfig{
		HijackMode: "replace",
		Contents:   []ScreensaverContent{},
		Source:     "屏保功能来源于SeewoMitM",
		EmitTime:   600,
	}
}

type ScreensaverContent struct {
	/* 所有都会有的 */
	// 类型
	Type string `json:"type,omitempty"`

	// 持续时间（单位：毫秒）
	Duration int `json:"duration,omitempty"`

	/* 绝大多数都有的 */
	// 路径
	Path string `json:"path,omitempty"`

	// 是否需要预加载
	RequirePreload bool `json:"requirePreload,omitempty"`

	/* 只有图片和视频有的 */
	// 适应方式
	Fit string `json:"fit,omitempty"`

	/* 只有视频有的 */
	// 是否静音
	Muted bool `json:"muted,omitempty"`

	/* 视频和Spine才有的*/
	// 是否循环播放
	Loop bool `json:"loop,omitempty"`

	/* 只有Spine有的 */
	// SpinePlayer的配置
	SpinePlayerConfig *SpinePlayerConfig `json:"spinePlayerConfig,omitempty"`
}

func NewScreensaverContent() *ScreensaverContent {
	return &ScreensaverContent{}
}

func NewScreensaverImageContent(path string, requirePreload bool, fit string, duration int) *ScreensaverContent {
	return &ScreensaverContent{
		Type:           "image",
		Path:           path,
		RequirePreload: requirePreload,
		Fit:            fit,
		Duration:       duration,
	}
}

func NewScreensaverVideoContent(path string, requirePreload bool, fit string, muted bool, loop bool, duration int) *ScreensaverContent {
	return &ScreensaverContent{
		Type:           "video",
		Path:           path,
		RequirePreload: requirePreload,
		Fit:            fit,
		Muted:          muted,
		Loop:           loop,
		Duration:       duration,
	}
}

func NewScreensaverSpineContent(path string, requirePreload bool, spinePlayerConfig *SpinePlayerConfig, duration int) *ScreensaverContent {
	return &ScreensaverContent{
		Type:              "spine",
		Path:              path,
		RequirePreload:    requirePreload,
		SpinePlayerConfig: spinePlayerConfig,
		Duration:          duration,
	}
}
