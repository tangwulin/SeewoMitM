package main

import (
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/request_handler"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	logFile := helper.GetLogFile("")
	log.InitGlobal(
		log.NewLogrusAdapt(&logrus.Logger{Out: io.MultiWriter(os.Stdout, logFile),
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks), Level: logrus.TraceLevel}))
}

func main() {
	upstreamPort, err := helper.GetUpstream(5)
	if err != nil {
		log.WithFields(log.Fields{"type": "GetUpstream"}).Error(err.Error())
		return
	}
	log.WithFields(log.Fields{"type": "Upstream"}).Info(fmt.Sprintf("upstream port:%d", upstreamPort))

	wsUrl := fmt.Sprintf("wss://127.0.0.1:%d/forward/SeewoHugoHttp/SeewoHugoService", upstreamPort)
	isWsAvailable := helper.TestWSAvailable(wsUrl, nil)
	if !isWsAvailable {
		log.WithFields(log.Fields{"type": "TestWSAvailable"}).Error("upstream websocket is not available")
		return
	}
	log.WithFields(log.Fields{"type": "TestWSAvailable"}).Info("upstream websocket is available")

	downstreamPort, err := helper.GetAvailablePort(14514)
	if err != nil {
		log.WithFields(log.Fields{"type": "GetAvailablePort"}).Error(err.Error())
		return
	}
	log.WithFields(log.Fields{"type": "Downstream"}).Info(fmt.Sprintf("downstream port:%d", downstreamPort))

	reqHandler := request_handler.RequestHandler(upstreamPort)

	currentPath, err := filepath.Abs(".")
	if err != nil {
		log.WithFields(log.Fields{"type": "GetCurrentPath"}).Error(err.Error())
		panic(err)
	}

	certPath := filepath.Join(currentPath, "server.crt")
	keyPath := filepath.Join(currentPath, "server.key")
	log.WithFields(log.Fields{"type": "TLS cert"}).Info(fmt.Sprintf("cert path:%s, key path:%s", certPath, keyPath))

	http.HandleFunc("/", reqHandler)

	log.WithFields(log.Fields{"type": "Server"}).Info(fmt.Sprintf("Listening on port %d", downstreamPort))

	err = http.ListenAndServeTLS(fmt.Sprintf(":%d", downstreamPort), certPath, keyPath, nil)
	//err = http.ListenAndServe(fmt.Sprintf(":%d", downstreamPort), nil)

	if err != nil {
		log.WithFields(log.Fields{"type": "Server"}).Error(err.Error())
	}
}
