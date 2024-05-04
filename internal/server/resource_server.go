package server

import (
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

var cacheMap = make(map[string]string)
var cacheDir string

func LaunchResourceServer(port int) error {
	err := BeforeLaunchResourceServer()
	if err != nil {
		log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Resource server start failed:%v", err.Error())
		return err
	}

	err = launchResourceServer(port)
	if err != nil {
		log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Resource server start failed:%v", err.Error())
		return err
	}
	return nil
}

func BeforeLaunchResourceServer() error {
	if helper.GetConfig().Cache == nil {
		return errors.New("cache config is nil")
	}

	if helper.GetConfig().Cache.CacheDir == "" {
		cacheDir = ".\\cache"
	} else {
		cacheDir = helper.GetConfig().Cache.CacheDir
	}

	// 检查缓存目录是否存在
	if !helper.PathExists(cacheDir) {
		err := os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Create cache dir failed:%v", err.Error())
			return err
		}
	}
	return nil
}

func launchResourceServer(port int) error {
	r := gin.New()
	r.GET("/getResource", func(c *gin.Context) {
		go func() {
			getResourceHandler(c)
		}()
	})

	err := r.Run(":" + strconv.Itoa(port))
	if err != nil {
		return err
	}
	return nil
}

func ReloadResourceServer() {
	err := BeforeLaunchResourceServer()
	if err != nil {
		log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Resource server reload failed:%v", err.Error())
		return
	}
}

func getResourceHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.HTML(400, "index.html", "url is empty")
		return
	}

	path, exist := cacheMap[url]

	// 如果map中不存在，检查缓存目录中是否存在
	if !exist {
		path = cacheDir + "\\" + helper.GetMD5(url)

		// 如果缓存目录中不存在，请求资源并缓存
		exist = helper.PathExists(path)
		if !exist {
			resp, err := http.Get(url)
			if err != nil {
				log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Get resource failed:%v", err.Error())
				c.HTML(400, "index.html", "Get resource failed")
				return
			}
			defer resp.Body.Close()

			cacheFile, err := os.Create(path)
			if err != nil {
				// 如果创建文件失败，直接返回资源
				log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Create cache file failed:%v", err.Error())
				c.Writer.WriteHeader(200)
				_, err = io.Copy(c.Writer, resp.Body)
				if err != nil {
					log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Copy response to client error:%v", err.Error())
				}
				c.Writer.Flush()
				return
			}
			defer cacheFile.Close()

			// 将资源写入缓存文件并返回给客户端
			multiWriter := io.MultiWriter(cacheFile, c.Writer)
			c.Writer.WriteHeader(200)
			_, err = io.Copy(multiWriter, resp.Body)
			if err != nil {
				log.WithFields(log.Fields{"type": "ResourceServer"}).Errorf("Copy response to client or writing to cache file error:%v", err.Error())
			}
			cacheMap[url] = path
			c.Writer.Flush()
			return
		}
		cacheMap[url] = path
		c.File(path)
		return
	}

	c.File(path)
	return
}
