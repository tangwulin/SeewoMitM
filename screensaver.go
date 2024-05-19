package main

import (
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"strings"
)

var screensaverContent DataContent

type DataContent struct {
	Mode           string        `json:"mode"`
	ImageList      []string      `json:"imageList"`
	ExtraPayload   *ExtraPayload `json:"extraPayload,omitempty"`
	Source         string        `json:"source"`
	Fit            string        `json:"fit"`
	PlayMode       string        `json:"playMode"`
	SwitchInterval int           `json:"switchInterval"`
	TextList       []TextItem    `json:"textList"`
}

func ParseScreensaverContent() DataContent {
	gc := GetConfig()
	if gc.ScreensaverConfig == nil {
		return DataContent{}
	}

	imgList := make([]string, 0)
	contents := make([]Content, 0)

	//TODO: 资源预加载还需要针对spine进行处理
	for _, c := range gc.ScreensaverConfig.Contents {
		switch c.Type {
		case "image":
			url := GetContentTrueUrl(c)

			imgList = append(imgList, url)
			contents = append(contents, *NewImageContent(url, c.Fit, c.Duration))
		case "video":
			contents = append(contents, *NewVideoContent(GetContentTrueUrl(c), c.Fit, c.Duration))
		case "spine":
			/*// 如果有RawDataURIs就直接使用
			if len(c.SpinePlayerConfig.RawDataURIs) > 0 {
				contents = append(contents, *NewSpineContent(c.Path, c.SpinePlayerConfig, c.Duration))
				continue
			}

			// 如果没有RawDataURIs就使用AtlasUrl和JsonUrl
			spinePlayerConfig := *c.SpinePlayerConfig

			if &spinePlayerConfig.AtlasUrl != nil && spinePlayerConfig.AtlasUrl != "" {
				spinePlayerConfig.AtlasUrl = GetResourceTrueUrl(spinePlayerConfig.AtlasUrl, c.RequirePreload)
			}

			if &spinePlayerConfig.JsonUrl != nil && spinePlayerConfig.JsonUrl != "" {
				spinePlayerConfig.JsonUrl = GetResourceTrueUrl(spinePlayerConfig.JsonUrl, c.RequirePreload)
			}

			if spinePlayerConfig.BackgroundImage != nil && spinePlayerConfig.BackgroundImage.Url != "" {
				spinePlayerConfig.BackgroundImage.Url = GetResourceTrueUrl(spinePlayerConfig.BackgroundImage.Url, c.RequirePreload)
			}*/

			contents = append(contents, *NewSpineContent(c.Path, c.SpinePlayerConfig, c.Duration, c.Scale, c.OffsetX, c.OffsetY))
		default:
			log.WithFields(log.Fields{"type": "ParseScreensaverContent"}).Error("unknown content type, content will be ignored:", c)
			continue
		}
	}

	textList := make([]TextItem, 0)
	// golang是没有鸭子类型吗？？？
	for _, t := range gc.ScreensaverConfig.TextList {
		textList = append(textList, TextItem{
			Content:    t.Content,
			Provenance: t.Provenance,
		})
	}

	return DataContent{
		Mode:           gc.ScreensaverConfig.Mode,
		ImageList:      imgList,
		ExtraPayload:   &ExtraPayload{ScreensaverContent: contents, ScreensaverSwitchInterval: gc.ScreensaverConfig.SwitchInterval},
		Source:         gc.ScreensaverConfig.Source,
		Fit:            gc.ScreensaverConfig.Fit,
		PlayMode:       gc.ScreensaverConfig.PlayMode,
		SwitchInterval: gc.ScreensaverConfig.SwitchInterval / 1000,
		TextList:       textList,
	}
}

func GetContentTrueUrl(c ScreensaverContent) string {
	if c.IsRequirePreload() {
		md5 := helper.MD5Sum([]byte(c.Path))
		PrepareResource(c.Path, md5)
		return GetResourceServerAddr() + "/" + md5
	}
	return c.Path
}

func GetResourceTrueUrl(url string, requirePreload bool) string {
	if strings.HasPrefix(url, "http") && requirePreload {
		return GetResourceServerAddr() + "/" + url
	}
	return url
}

func LoadScreensaverContent() {
	screensaverContent = ParseScreensaverContent()
}

func GetScreensaverContent() *DataContent {
	return &screensaverContent
}
