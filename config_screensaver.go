package main

import (
	"SeewoMitM/model"
	"errors"
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

func GetScreensaverContentList() []model.ScreensaverContent {
	return globalConfig.ScreensaverConfig.Contents
}

func GetScreensaverContentByID(id int) (result model.ScreensaverContent, err error) {
	for _, v := range globalConfig.ScreensaverConfig.Contents {
		if v.ID == id {
			return v, nil
		}
	}

	return result, errors.New("未找到该内容")
}

func AddScreensaverContent(content model.ScreensaverContent) (id int, err error) {
	contentsLen := len(globalConfig.ScreensaverConfig.Contents)
	if contentsLen > 0 {
		content.ID = globalConfig.ScreensaverConfig.Contents[contentsLen-1].ID + 1
	} else {
		content.ID = 0
	}

	globalConfig.ScreensaverConfig.Contents = append(globalConfig.ScreensaverConfig.Contents, content)
	err = SaveConfig()
	if err != nil {
		return -1, err
	}
	return content.ID, nil
}

func DeleteScreensaverContent(id int) (err error) {
	for i, content := range globalConfig.ScreensaverConfig.Contents {
		if content.ID == id {
			globalConfig.ScreensaverConfig.Contents = append(globalConfig.ScreensaverConfig.Contents[:i], globalConfig.ScreensaverConfig.Contents[i+1:]...)
			break
		}
		return errors.New("未找到该内容")
	}

	err = SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func UpdateScreensaverContent(id int, newContent model.ScreensaverContent) (err error) {
	for i, content := range globalConfig.ScreensaverConfig.Contents {
		if content.ID == id {
			globalConfig.ScreensaverConfig.Contents[i] = newContent
			break
		}
		return errors.New("未找到该内容")
	}

	err = SaveConfig()
	if err != nil {
		return err
	}
	return nil
}
