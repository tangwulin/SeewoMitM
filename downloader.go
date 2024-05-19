package main

import (
	"SeewoMitM/internal/log"
	"io"
	"net/http"
	"os"
)

type downloadTask struct {
	Url  string
	Path string
}

var downloadQueue = make(chan downloadTask, 1000)

func LaunchDownloader(threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			log.WithFields(log.Fields{"type": "Downloader"}).Infof("downloader %d# started", i)
			for {
				task := <-downloadQueue
				file, err := os.OpenFile(task.Path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
				if err != nil {
					log.WithFields(log.Fields{"type": "Downloader"}).Error("failed to open file:", task.Path)
					continue
				}
				defer file.Close()

				resp, err := http.Get(task.Url)
				if err != nil {
					log.WithFields(log.Fields{"type": "Downloader"}).Error("failed to download file:", task.Url)
					continue
				}
				defer resp.Body.Close()

				_, err = io.Copy(file, resp.Body)
				if err != nil {
					log.WithFields(log.Fields{"type": "Downloader"}).Error("failed to write file:", task.Path)
					continue
				}
				log.WithFields(log.Fields{"type": "Downloader"}).Info("downloaded file:", task.Path)
			}
		}()
	}
}

func AddDownloadTask(url, path string) {
	downloadQueue <- downloadTask{Url: url, Path: path}
}
