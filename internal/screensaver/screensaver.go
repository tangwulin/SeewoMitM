package screensaver

import (
	"SeewoMitM/internal/config"
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/resource"
	"strings"
)

var screensaverDataContents DataContent

type DataContent struct {
	ImageList    []string     `json:"imageList"`
	ExtraPayload ExtraPayload `json:"extraPayload,omitempty"`
}

type ExtraPayload struct {
	ScreensaverContent []Content `json:"screensaverContent"`
}

func ParseScreensaverContent() DataContent {
	gc := helper.GetConfig()
	if gc.ScreensaverConfig == nil {
		return DataContent{}
	}

	imgList := make([]string, 0)
	contents := make([]Content, 0)

	for _, c := range gc.ScreensaverConfig.Contents {
		switch c.Type {
		case "image":
			url := GetContentTrueUrl(c)

			imgList = append(imgList, url)
			contents = append(contents, *NewImageContent(url, c.Fit, c.Duration))
		case "video":
			contents = append(contents, *NewVideoContent(GetContentTrueUrl(c), c.Fit, c.Duration))
		case "spine":
			// 如果有RawDataURIs就直接使用
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
			}

			contents = append(contents, *NewSpineContent(c.Path, &spinePlayerConfig, c.Duration))
		default:
			log.WithFields(log.Fields{"type": "ParseScreensaverContent"}).Error("unknown content type, content will be ignored:", c)
			continue
		}
	}

	return DataContent{
		ImageList:    imgList,
		ExtraPayload: ExtraPayload{ScreensaverContent: contents},
	}
}

func GetContentTrueUrl(c config.ScreensaverContent) string {
	if c.IsRequirePreload() {
		md5 := helper.MD5Sum([]byte(c.Path))
		resource.PrepareResource(c.Path, md5)
		return resource.GetResourceServerAddr() + "/" + md5
	}
	return c.Path
}

func GetResourceTrueUrl(url string, requirePreload bool) string {
	if strings.HasPrefix(url, "http") && requirePreload {
		return resource.GetResourceServerAddr() + "/" + url
	}
	return url
}

func LoadScreensaverData() {
	screensaverDataContents = ParseScreensaverContent()
}

func GetScreensaverData() *DataContent {
	return &screensaverDataContents
}
