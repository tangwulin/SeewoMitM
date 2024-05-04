package services

import (
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/request_handler"
	"crypto/tls"
	"embed"
	"fmt"
	"net/http"
	"strconv"
)

func LaunchMitMService(downstreamPort int, upstreamPort int, certFiles embed.FS) error {
	// 读取证书文件
	certContent, err := certFiles.ReadFile("server.crt")
	if err != nil {
		log.WithFields(log.Fields{"type": "ReadCertFile"}).Error(err.Error())
		return err
	}
	keyContent, err := certFiles.ReadFile("server.key")
	if err != nil {
		log.WithFields(log.Fields{"type": "ReadKeyFile"}).Error(err.Error())
		return err
	}
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
	return err
}
