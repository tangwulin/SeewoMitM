package main

import "SeewoMitM/model"

func NewImageContent(path string, fit string, duration int) *model.ScreensaverContentPayload {
	return &model.ScreensaverContentPayload{
		Type:     "image",
		Path:     path,
		Fit:      fit,
		Duration: duration,
	}
}

func NewVideoContent(path string, fit string, duration int) *model.ScreensaverContentPayload {
	return &model.ScreensaverContentPayload{
		Type:     "video",
		Path:     path,
		Fit:      fit,
		Duration: duration,
	}
}

func NewSpineContent(path string, spinePlayerConfig *model.SpinePlayerConfig, duration int, scale, offsetX, offsetY float64) *model.ScreensaverContentPayload {
	return &model.ScreensaverContentPayload{
		Type:              "spine",
		Path:              path,
		SpinePlayerConfig: spinePlayerConfig,
		Duration:          duration,
		Scale:             scale,
		OffsetX:           offsetX,
		OffsetY:           offsetY,
	}
}
