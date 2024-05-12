package main

type Content struct {
	Type              string             `json:"type,omitempty"`
	Path              string             `json:"path,omitempty"`
	Fit               string             `json:"fit,omitempty"`
	SpinePlayerConfig *SpinePlayerConfig `json:"spinePlayerConfig,omitempty"`
	Duration          int                `json:"duration,omitempty"`
	Scale             float64            `json:"scale,omitempty"`
	OffsetX           float64            `json:"offsetX,omitempty"`
	OffsetY           float64            `json:"offsetY,omitempty"`
}

func NewImageContent(path string, fit string, duration int) *Content {
	return &Content{
		Type:     "image",
		Path:     path,
		Fit:      fit,
		Duration: duration,
	}
}

func NewVideoContent(path string, fit string, duration int) *Content {
	return &Content{
		Type:     "video",
		Path:     path,
		Fit:      fit,
		Duration: duration,
	}
}

func NewSpineContent(path string, spinePlayerConfig *SpinePlayerConfig, duration int, scale, offsetX, offsetY float64) *Content {
	return &Content{
		Type:              "spine",
		Path:              path,
		SpinePlayerConfig: spinePlayerConfig,
		Duration:          duration,
		Scale:             scale,
		OffsetX:           offsetX,
		OffsetY:           offsetY,
	}
}
