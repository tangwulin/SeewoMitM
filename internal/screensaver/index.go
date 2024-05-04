package screensaver

import (
	"SeewoMitM/internal/config"
	"SeewoMitM/internal/downloader"
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/spine"
	"path"
	"strings"
)

var screenSaverContent []ScreenSaverContent

type ScreenSaverContent struct {
	Type         string        `json:"type"`
	Url          string        `json:"url"`
	Fit          string        `json:"fit"`
	SpineOptions spine.Options `json:"spineOptions"`
}

func GetScreenSaverContent() []ScreenSaverContent {
	return screenSaverContent
}

func LoadScreenSaverConfig(c *config.Config) {
	if c == nil || len(c.ScreenSaverHijackContent) == 0 {
		return
	}

	for _, v := range c.ScreenSaverHijackContent {
		switch v.Type {
		case "html":
			//TODO
			continue
		case "image":
			// 不是本地文件才需要预加载
			if v.RequirePreload && strings.HasPrefix(v.Path, "http") {
				filePath := path.Join(config.GetResourceServerAddr(), helper.GetMD5(v.Path))
				downloader.AddDownloadTask(v.Path, filePath)
				screenSaverContent = append(screenSaverContent, ScreenSaverContent{
					Type: "image",
					Url:  filePath,
				})
				log.WithFields(log.Fields{"type": "ScreenSaver"}).Infof("Add preload image:%v", filePath)
			} else {
				screenSaverContent = append(screenSaverContent, ScreenSaverContent{
					Type: "image",
					Url:  v.Path,
				})
			}
		case "video":
			// 不是本地文件才需要预加载
			if v.RequirePreload && strings.HasPrefix(v.Path, "http") {
				filePath := path.Join(config.GetResourceServerAddr(), helper.GetMD5(v.Path))
				downloader.AddDownloadTask(v.Path, filePath)
				screenSaverContent = append(screenSaverContent, ScreenSaverContent{
					Type: "video",
					Url:  filePath,
				})
			} else {
				screenSaverContent = append(screenSaverContent, ScreenSaverContent{
					Type: "video",
					Url:  v.Path,
				})
			}
		case "spine":
			if v.SpineConfig == nil {
				continue
			}
			spineOptions := v.SpineConfig
			if v.RequirePreload {
				if strings.HasPrefix(v.SpineConfig.JsonUrl, "http") {
					filePath := path.Join(config.GetResourceServerAddr(), helper.GetMD5(v.SpineConfig.JsonUrl))
					downloader.AddDownloadTask(v.SpineConfig.JsonUrl, filePath)
					spineOptions.JsonUrl = filePath
				}
				if strings.HasPrefix(v.SpineConfig.SkelUrl, "http") {
					filePath := path.Join(config.GetResourceServerAddr(), helper.GetMD5(v.SpineConfig.AtlasUrl))
					downloader.AddDownloadTask(v.SpineConfig.AtlasUrl, filePath)
					spineOptions.AtlasUrl = filePath
				}
				if strings.HasPrefix(v.SpineConfig.AtlasUrl, "http") {
					filePath := path.Join(config.GetResourceServerAddr(), helper.GetMD5(v.SpineConfig.AtlasUrl))
					downloader.AddDownloadTask(v.SpineConfig.AtlasUrl, filePath)
					spineOptions.AtlasUrl = filePath
				}
			}

			screenSaverContent = append(screenSaverContent, ScreenSaverContent{
				Type:         "spine",
				SpineOptions: *spineOptions,
			})
		//case "imageDirectory":
		// TODO
		//	watcher.AddWatch(v.Path)
		default:
			continue
		}
	}
	//err := watcher.ReloadWatchList()
	//if err != nil {
	//	return
	//}
}
