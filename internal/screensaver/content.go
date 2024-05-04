package screensaver

import (
	"SeewoMitM/internal/spine"
)

type Content struct {
	Type              string                   `json:"type,omitempty"`
	Path              string                   `json:"path,omitempty"`
	Fit               string                   `json:"fit,omitempty"`
	SpinePlayerConfig *spine.SpinePlayerConfig `json:"spinePlayerConfig,omitempty"`
	Duration          int                      `json:"duration,omitempty"`
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

func NewSpineContent(path string, spinePlayerConfig *spine.SpinePlayerConfig, duration int) *Content {
	return &Content{
		Type:              "spine",
		Path:              path,
		SpinePlayerConfig: spinePlayerConfig,
		Duration:          duration,
	}
}
