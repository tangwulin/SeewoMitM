package main

import (
	"SeewoMitM/internal/log"
	"SeewoMitM/model"
)

var screensaverContent model.DataContent

func ParseScreensaverContent() model.DataContent {
	screensaverConfig := globalConfig.ScreensaverConfig

	imgList := make([]string, 0, len(screensaverConfig.Contents)/2)
	contents := make([]model.ScreensaverContentPayload, 0, len(screensaverConfig.Contents))

	//TODO: 资源预加载还需要针对spine进行处理
	for _, c := range screensaverConfig.Contents {
		switch c.Type {
		case "image":
			url := c.Path

			imgList = append(imgList, url)
			contents = append(contents, *NewImageContent(url, c.Fit, c.Duration))
		case "video":
			contents = append(contents, *NewVideoContent(c.Path, c.Fit, c.Duration))
		case "spine":
			contents = append(contents, *NewSpineContent(c.Path, c.SpinePlayerConfig, c.Duration, c.Scale, c.OffsetX, c.OffsetY))
		default:
			log.WithFields(log.Fields{"type": "ParseScreensaverContent"}).Error("unknown content type, content will be ignored:", c)
			continue
		}
	}

	textList := make([]model.TextItem, 0)
	// golang是没有鸭子类型吗？？？
	for _, t := range screensaverConfig.TextList {
		textList = append(textList, model.TextItem{
			Content:    t.Content,
			Provenance: t.Provenance,
		})
	}

	return model.DataContent{
		Mode:           screensaverConfig.Mode,
		ImageList:      imgList,
		ExtraPayload:   &model.ExtraPayload{ScreensaverContent: contents, ScreensaverSwitchInterval: screensaverConfig.SwitchInterval},
		Source:         screensaverConfig.Source,
		Fit:            screensaverConfig.Fit,
		PlayMode:       screensaverConfig.PlayMode,
		SwitchInterval: screensaverConfig.SwitchInterval / 1000,
		TextList:       textList,
	}
}

func LoadScreensaverContent() {
	screensaverContent = ParseScreensaverContent()
}

func GetScreensaverContent() *model.DataContent {
	return &screensaverContent
}
