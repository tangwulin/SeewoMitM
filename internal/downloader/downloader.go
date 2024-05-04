package downloader

import (
	"SeewoMitM/internal/log"
	"io"
	"net/http"
	"os"
	"path"
)

type DownloadTask struct {
	Url  string
	Path string
}

var downloadQueue = make(chan DownloadTask, 1000)

func AddDownloadTask(url string, path string) {
	downloadQueue <- DownloadTask{Url: url, Path: path}
	log.WithFields(log.Fields{"type": "Downloader"}).Infof("Add download task, url:%v, path:%v", url, path)
}

func StartDownloader() {
	for {
		task := <-downloadQueue
		go func() {
			downloadDir := path.Dir(task.Path)
			// 创建下载目录
			err := os.MkdirAll(downloadDir, os.ModePerm)
			if err != nil {
				log.WithFields(log.Fields{"type": "Downloader"}).Errorf("Create download dir failed:%v", err.Error())
				return
			}
			// 下载文件
			file, err := os.OpenFile(task.Path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				log.WithFields(log.Fields{"type": "Downloader"}).Errorf("Open file failed:%v", err.Error())
				return
			}
			defer file.Close()
			resp, err := http.Get(task.Url)
			if err != nil {
				log.WithFields(log.Fields{"type": "Downloader"}).Errorf("http.Get failed:%v", err.Error())
				return
			}
			defer resp.Body.Close()
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				log.WithFields(log.Fields{"type": "Downloader"}).Errorf("io.Copy failed:%v", err.Error())
				return
			}
			log.WithFields(log.Fields{"type": "Downloader"}).Infof("Download file success, url:%v, file: %v", task.Url, task.Path)
		}()
	}
}
