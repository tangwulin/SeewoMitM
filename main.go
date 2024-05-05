package main

import (
	"SeewoMitM/internal/config"
	"SeewoMitM/internal/connection"
	"SeewoMitM/internal/downloader"
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/manage"
	"SeewoMitM/internal/mitm"
	"SeewoMitM/internal/resource"
	"SeewoMitM/internal/screensaver"
	"SeewoMitM/internal/timer"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

//go:embed server.crt server.key
var certFiles embed.FS

func main() {
	// 读取命令行参数
	configFilePathPtr := flag.String("config", "", "配置文件路径")
	logFilePathPtr := flag.String("log", "", "日志文件路径")
	upstreamPortPtr := flag.Int("upstream", 0, "上游端口")
	downstreamPortPtr := flag.Int("downstream", 0, "下游端口")
	logLevelPtr := flag.String("logLevel", "", "日志级别")
	//runAsDaemonPtr := flag.Bool("daemon", false, "是否以守护进程运行")

	// 解析命令行参数
	flag.Parse()

	var configFilePath = "./config.json"
	var logDir = ""
	var upstreamPort int
	var downstreamPort int
	var logLevel string

	if *configFilePathPtr != "" {
		configFilePath = *configFilePathPtr
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if configFilePath == "./config.json" {
			fmt.Println("config file not found, creating default config file...")

			// 创建默认配置文件
			file, err := os.Create(configFilePath)
			if err != nil {
				fmt.Printf("create config file error: %v\n", err)
				return
			}
			defer file.Close()

			defaultConfig := config.Config{
				LogLevel:          "info",
				ScreensaverConfig: config.NewScreensaverConfig(),
			}

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")
			err = encoder.Encode(defaultConfig)
			if err != nil {
				fmt.Printf("write config file error: %v\n", err)
				return
			}
			file.Close()
		} else {
			fmt.Printf("config file not found\n")
			return
		}
	}

	//检查是否有配置文件
	configs, err := helper.ReadAndParseConfig(configFilePath)
	if err != nil {
		fmt.Printf("ReadAndParseConfig error: %v\n", err)
		panic(err)
		return
	}

	helper.SetConfig(*configs)

	// 检测有没有指定日志文件路径
	if *logFilePathPtr != "" {
		// 检查日志文件路径是否存在
		if _, err := os.Stat(*logFilePathPtr); os.IsNotExist(err) {
			// 创建日志文件夹
			err = os.MkdirAll(*logFilePathPtr, os.ModePerm)
			if err != nil {
				panic(err)
				return
			}
		}
		logDir = *logFilePathPtr
	}

	logLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}

	// 检测有没有指定日志级别
	if *logLevelPtr != "" {
		for _, v := range logLevels {
			if *logLevelPtr == v {
				logLevel = *logLevelPtr
				break
			}
		}
		if logLevel == "" {
			logLevel = configs.LogLevel
		}
	}

	// 初始化日志
	logFile := helper.GetLogFile(logDir)
	log.InitGlobal(
		log.NewLogrusAdapt(&logrus.Logger{Out: io.MultiWriter(os.Stdout, logFile),
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks), Level: logrus.Level(log.FindLevel(logLevel))}))

	// 检测上游端口是否指定，如果没有指定则自动获取
	if *upstreamPortPtr != 0 {
		upstreamPort = *upstreamPortPtr
	} else {
		upstreamPort, err = helper.GetUpstream(5)
		if err != nil {
			log.WithFields(log.Fields{"type": "GetUpstream"}).Error(err.Error())
			return
		}
	}
	log.WithFields(log.Fields{"type": "Upstream"}).Info(fmt.Sprintf("upstream port:%d", upstreamPort))

	// 检测上游Websocket是否可用
	wsUrl := fmt.Sprintf("wss://127.0.0.1:%d/forward/SeewoHugoHttp/SeewoHugoService", upstreamPort)
	isWsAvailable := helper.TestWSAvailable(wsUrl, nil)
	if !isWsAvailable {
		log.WithFields(log.Fields{"type": "TestWSAvailable"}).Error("upstream websocket is not available")
		return
	}
	log.WithFields(log.Fields{"type": "TestWSAvailable"}).Info("upstream websocket is available")

	// 检测下游端口是否指定，如果没有指定则自动获取
	if *downstreamPortPtr != 0 {
		downstreamPort = *downstreamPortPtr
	} else {
		downstreamPort, err = helper.GetAvailablePort(14514)
		if err != nil {
			log.WithFields(log.Fields{"type": "GetAvailablePort"}).Error(err.Error())
			return
		}
	}
	log.WithFields(log.Fields{"type": "Downstream"}).Info(fmt.Sprintf("downstream port:%d", downstreamPort))

	downloader.LaunchDownloader(2)
	screensaver.LoadScreensaverContent()
	resource.LaunchResourceService(14515, "./resource")

	wg := sync.WaitGroup{}
	wg.Add(2)

	// 启动管理服务
	go func() {
		err = manage.LaunchManageServer(14516)
		if err != nil {
			log.WithFields(log.Fields{"type": "LaunchManageServer"}).Error(err.Error())
		}
		wg.Done()
	}()

	// 启动服务端
	go func() {
		err = mitm.LaunchMitMService(downstreamPort, upstreamPort, certFiles)
		if err != nil {
			log.WithFields(log.Fields{"type": "LaunchMitMService"}).Error(err.Error())
		}
		wg.Done()
	}()

	err = helper.WriteMitMPortToRegistry(downstreamPort)
	if err != nil {
		log.WithFields(log.Fields{"type": "WriteMitMPortToRegistry"}).Error(fmt.Sprintf("Count not to write MitM port to registry: %v", err))
	} else {
		log.WithFields(log.Fields{"type": "WriteMitMPortToRegistry"}).Info(fmt.Sprintf("MitM port is written to registry: %d", downstreamPort))
	}

	err = helper.KillProcessByName("SeewoServiceAssistant.exe")
	if err != nil {
		log.WithFields(log.Fields{"type": "RelaunchSeewoServiceAssistant"}).Error(fmt.Sprintf("Count not to kill SeewoServiceAssistant.exe: %v", err))
	} else {
		log.WithFields(log.Fields{"type": "RelaunchSeewoServiceAssistant"}).Info("SeewoServiceAssistant.exe is relaunched")
	}

	//启动屏保定时器
	go func() {
		err := timer.LaunchScreenSaverTimer(time.Duration(configs.ScreensaverConfig.EmitTime)*time.Second, func() {
			cp := *connection.GetConnectionPool()
			for _, v := range cp {
				if v.URL == "/forward/SeewoHugoHttp/SeewoHugoService" {
					err := v.DownstreamConn.WriteMessage(websocket.TextMessage, []byte(`{
    "data": {
        "imageList": [
            "D:\\85499466.jpg",
            "D:\\532421.jpg",
            "D:\\650142.jpg",
            "D:\\124177.jpg",
            "D:\\1325365.jpg"
        ],
        "materialSource": "屏保功能来源于中国人口吧",
        "extraPayload": [
            "屏保功能来源于中国人口吧",
            {
                "type": "image",
                "path": "D:\\85499466.jpg"
            }
        ],
        "pictureSizeType": 1,
        "playMode": 0,
        "switchInterval": 5,
        "textList": [
            {
                "content": "aaa",
                "provenance": "bbb"
            },
            {
                "content": "foo",
                "provenance": "bar"
            }
        ]
    },
    "traceId": "0C89A601-B51D-488D-87EC-5862CE75ABE7",
    "url": "/displayScreenSaver"
}`))
					if err != nil {
						return
					}
				}
			}
		})
		if err != nil {
			return
		}
	}()

	wg.Wait()
}
