package model

type ScreensaverContentPayload struct {
	Type              string             `json:"type,omitempty"`
	Path              string             `json:"path,omitempty"`
	Fit               string             `json:"fit,omitempty"`
	SpinePlayerConfig *SpinePlayerConfig `json:"spinePlayerConfig,omitempty"`
	Duration          int                `json:"duration,omitempty"`
	Scale             float64            `json:"scale,omitempty"`
	OffsetX           float64            `json:"offsetX,omitempty"`
	OffsetY           float64            `json:"offsetY,omitempty"`
}
