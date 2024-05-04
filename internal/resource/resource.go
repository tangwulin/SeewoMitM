package resource

import (
	"SeewoMitM/internal/downloader"
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"os"
	"path"
)

var resourceDir string
var resourceServerAddr string

func SetResourceDir(path string) error {
	exists, _ := helper.PathExists(path)
	if !exists {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.WithFields(log.Fields{"type": "Resource", "path": path}).Error("failed to create resource directory")
			return err
		}
		log.WithFields(log.Fields{"type": "Resource", "path": path}).Info("resource directory created")
	}
	resourceDir = path
	return nil
}

func PrepareResource(url, filename string) {
	if filename == "" {
		filename = helper.MD5Sum([]byte(url))
	}
	filePath := path.Join(resourceDir, filename)
	exists, _ := helper.PathExists(filePath)
	if exists {
		fileInfo, err := helper.GetFileInfo(filePath)
		if err != nil {
			log.WithFields(log.Fields{"type": "Resource"}).Error("failed to get file info")
			return
		}
		if fileInfo.Size() > 0 {
			return
		}
		err = os.Remove(filePath)
		if err != nil {
			return
		}
	}
	downloader.AddDownloadTask(url, path.Join(resourceDir, filename))
}

func LaunchResourceService(port int, path string) {
	err := SetResourceDir(path)
	if err != nil {
		log.WithFields(log.Fields{"type": "Resource"}).Error("failed to set resource directory")
		return
	}
}

func GetResourceServerAddr() string {
	return resourceServerAddr
}
