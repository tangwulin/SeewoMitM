package model

type ScreensaverConfig struct {
	// 劫持模式
	Mode string `json:"mode"`

	// 劫持内容
	Contents []ScreensaverContent `json:"contents"`

	// 右下角来源显示
	Source string `json:"source"`

	// 触发时间
	EmitTime int `json:"emitTime" default:"600"`

	// 内容适应模式
	Fit string `json:"fit"`

	// 轮播模式
	PlayMode string `json:"playMode"`

	// 轮播间隔（单位：毫秒）
	SwitchInterval int `json:"switchInterval"`

	// 文字列表
	TextList []struct {
		Content    string `json:"content"`
		Provenance string `json:"provenance"`
	} `json:"textList"`
}
