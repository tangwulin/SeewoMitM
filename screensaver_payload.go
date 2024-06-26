package main

import (
	"github.com/google/uuid"
	"strings"
)

type ScreensaverPayload struct {
	Data    Data   `json:"data"`
	TraceId string `json:"traceId"`
	Url     string `json:"url"`
}

type Data struct {
	ImageList       []string      `json:"imageList"`
	MaterialSource  string        `json:"materialSource"`
	ExtraPayload    *ExtraPayload `json:"extraPayload,omitempty"`
	PictureSizeType int           `json:"pictureSizeType"`
	PlayMode        int           `json:"playMode"`
	SwitchInterval  int           `json:"switchInterval"`
	TextList        []TextItem    `json:"textList"`
}

type ExtraPayload struct {
	ScreensaverSwitchInterval int       `json:"screensaverSwitchInterval"`
	ScreensaverContent        []Content `json:"screensaverContent"`
}

type TextItem struct {
	Content    string `json:"content"`
	Provenance string `json:"provenance"`
}

func GetScreensaverPayload() *ScreensaverPayload {
	content := GetScreensaverContent()

	var pictureSizeType int
	switch content.Fit {
	case "contain":
		pictureSizeType = 0
	case "cover":
		pictureSizeType = 1
	}

	return &ScreensaverPayload{
		Data: Data{
			ImageList:       content.ImageList,
			MaterialSource:  content.Source,
			ExtraPayload:    content.ExtraPayload,
			PictureSizeType: pictureSizeType,
			PlayMode:        0,
			SwitchInterval:  content.SwitchInterval,
			TextList:        content.TextList,
		},
		TraceId: strings.ToUpper(uuid.New().String()),
		Url:     "/displayScreenSaver",
	}
}
