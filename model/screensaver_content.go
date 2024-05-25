package model

import "strings"

type ScreensaverContent struct {
	/* 所有都会有的 */

	// ID
	ID int `json:"id,omitempty"`

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
	Scale             float64            `json:"scale,omitempty"`
	OffsetX           float64            `json:"offsetX,omitempty"`
	OffsetY           float64            `json:"offsetY,omitempty"`
}

func (content ScreensaverContent) IsRequirePreload() bool {
	return content.RequirePreload && strings.HasPrefix(content.Path, "http")
}
