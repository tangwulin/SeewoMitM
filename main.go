package main

import (
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
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

// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
func main() {
	// 读取命令行参数
	configFilePathPtr := flag.String("config", "", "配置文件路径")
	logFilePathPtr := flag.String("log", "", "日志文件路径")
	upstreamPortPtr := flag.Int("upstream", 0, "上游端口")
	downstreamPortPtr := flag.Int("downstream", 0, "下游端口")
	logLevelPtr := flag.String("logLevel", "", "日志级别")

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

			defaultConfig := Config{
				LogLevel:          "info",
				ScreensaverConfig: NewScreensaverConfig(),
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
	configs, err := ReadAndParseConfig(configFilePath)
	if err != nil {
		fmt.Printf("ReadAndParseConfig error: %v\n", err)
		panic(err)
		return
	}

	SetConfig(*configs)

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

	LaunchDownloader(2)
	LoadScreensaverContent()
	LaunchResourceService(14515, "./resource")

	wg := sync.WaitGroup{}
	wg.Add(2)

	// 启动管理服务
	go func() {
		managePort, err := helper.GetAvailablePort(11451)
		if err != nil {
			log.WithFields(log.Fields{"type": "GetAvailablePort"}).Error(err.Error())
			log.WithFields(log.Fields{"type": "Manage"}).Error("could not get available manage port")
			wg.Done()
			return
		}
		log.WithFields(log.Fields{"type": "Manage"}).Info(fmt.Sprintf("manage port:%d", managePort))
		err = LaunchManageServer(managePort)
		if err != nil {
			log.WithFields(log.Fields{"type": "LaunchManageServer"}).Error(err.Error())
		}
		wg.Done()
	}()

	// 启动服务端
	go func() {
		err = LaunchMitMService(downstreamPort, upstreamPort, certFiles)
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
		err := LaunchScreenSaverTimer(time.Duration(configs.ScreensaverConfig.EmitTime)*time.Second, func() {
			cp := *GetConnectionPool()
			for _, v := range cp {
				if v.URL == "/forward/SeewoHugoHttp/SeewoHugoService" {
					payload := GenScreensaverPayload()
					jsonData, err := json.Marshal(payload)
					if err != nil {
						log.WithFields(log.Fields{"type": "GenScreensaverPayload"}).Error(err.Error())
						continue
					}
					err = v.DownstreamConn.WriteMessage(websocket.TextMessage, jsonData)
					if err != nil {
						log.WithFields(log.Fields{"type": "SendPayload"}).Error(err.Error())
						continue
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
