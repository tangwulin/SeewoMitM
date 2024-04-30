package main

import (
	"SeewoMitM/internal/config"
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/request_handler"
	"crypto/tls"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
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
				LogLevel:                 "info",
				ScreenSaverHijackMode:    config.ScreenSaverHijackModeReplaceAll,
				ScreenSaverHijackContent: make([]config.ScreenSaverHijackContent, 0),
				ScreenSaverEmitTime:      600}

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")
			err = encoder.Encode(defaultConfig)
			if err != nil {
				fmt.Printf("write config file error: %v\n", err)
				return
			}
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

	// 读取证书文件
	certContent, err := certFiles.ReadFile("server.crt")
	if err != nil {
		log.WithFields(log.Fields{"type": "ReadCertFile"}).Error(err.Error())
		return
	}
	keyContent, err := certFiles.ReadFile("server.key")
	if err != nil {
		log.WithFields(log.Fields{"type": "ReadKeyFile"}).Error(err.Error())
		return
	}

	// 启动服务端
	log.WithFields(log.Fields{"type": "Server"}).Info(fmt.Sprintf("Listening on port %d", downstreamPort))

	cert, _ := tls.X509KeyPair(certContent, keyContent)

	s := &http.Server{
		Addr:      ":" + strconv.Itoa(downstreamPort),
		TLSConfig: &tls.Config{Certificates: append([]tls.Certificate{}, cert)},
	}

	reqHandler := request_handler.RequestHandler(upstreamPort)
	mux := http.NewServeMux()
	mux.HandleFunc("/", reqHandler)
	s.Handler = mux
	err = s.ListenAndServeTLS("", "")

	if err != nil {
		log.WithFields(log.Fields{"type": "Server"}).Error(err.Error())
	}
}
