package main

import (
	"SeewoMitM/model"
	"github.com/google/uuid"
	"strings"
)

func GenScreensaverPayload() *model.ScreensaverPayload {
	content := GetScreensaverContent()

	var pictureSizeType int
	switch content.Fit {
	case "contain":
		pictureSizeType = 0
	case "cover":
		pictureSizeType = 1
	}

	return &model.ScreensaverPayload{
		Data: model.ScreensaverPayloadData{
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
