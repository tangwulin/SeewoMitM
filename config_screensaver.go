package main

import (
	"SeewoMitM/model"
)

func NewScreensaverConfig() *model.ScreensaverConfig {
	return &model.ScreensaverConfig{
		Mode:           "replace",
		Contents:       []model.ScreensaverContent{},
		Source:         "屏保功能来源于SeewoMitM",
		EmitTime:       600,
		Fit:            "contain",
		PlayMode:       "sequence",
		SwitchInterval: 5000,
		TextList: []struct {
			Content    string `json:"content"`
			Provenance string `json:"provenance"`
		}{},
	}
}

func NewScreensaverContent() *model.ScreensaverContent {
	return &model.ScreensaverContent{}
}

func NewScreensaverImageContent(path string, requirePreload bool, fit string, duration int) *model.ScreensaverContent {
	return &model.ScreensaverContent{
		Type:           "image",
		Path:           path,
		RequirePreload: requirePreload,
		Fit:            fit,
		Duration:       duration,
	}
}

func NewScreensaverVideoContent(path string, requirePreload bool, fit string, muted bool, loop bool, duration int) *model.ScreensaverContent {
	return &model.ScreensaverContent{
		Type:           "video",
		Path:           path,
		RequirePreload: requirePreload,
		Fit:            fit,
		Muted:          muted,
		Loop:           loop,
		Duration:       duration,
	}
}

func NewScreensaverSpineContent(spinePlayerConfig *model.SpinePlayerConfig, requirePreload bool, duration int) *model.ScreensaverContent {
	return &model.ScreensaverContent{
		Type:              "spine",
		RequirePreload:    requirePreload,
		SpinePlayerConfig: spinePlayerConfig,
		Duration:          duration,
		Scale:             1,
		OffsetX:           0,
		OffsetY:           0,
	}
}
