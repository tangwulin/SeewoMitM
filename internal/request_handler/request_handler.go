package request_handler

import (
	"SeewoMitM/internal/log"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RequestHandler(upstreamPort int) func(w http.ResponseWriter, r *http.Request) {
	httpsUpstreamUrl, _ := url.Parse(fmt.Sprintf("https://localhost:%d", upstreamPort))
	forwardProxy := httputil.NewSingleHostReverseProxy(httpsUpstreamUrl)

	type wsMessage struct {
		messageType int
		payload     []byte
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r) {
			if log.IsLevelEnabled(log.TraceLevel) {
				bodyReader, err := r.GetBody()
				if err != nil {
					log.WithFields(log.Fields{"type": "Forward_HTTP_" + r.Method, "url": httpsUpstreamUrl.Path + r.RequestURI}).Error(err.Error())
				} else {
					body, _ := io.ReadAll(bodyReader)
					log.WithFields(log.Fields{"type": "Forward_HTTP_" + r.Method, "url": httpsUpstreamUrl.Path + r.RequestURI}).Trace(string(body))
				}
			}
			forwardProxy.ServeHTTP(w, r)
			return
		}

		upload := make(chan wsMessage, 20)
		download := make(chan wsMessage, 20)
		errChan := make(chan error, 20)

		downstream, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.WithFields(log.Fields{"type": "WS_Downstream_Upgrade"}).Error(err.Error())
			return
		} else {
			log.WithFields(log.Fields{"type": "WS_Downstream_Upgrade"}).Info(fmt.Sprintf("Downstream Websocket upgrade success, url:%s", r.RequestURI))
		}

		wssUpstreamUrl := fmt.Sprintf("wss://localhost:%d%s", upstreamPort, r.RequestURI)

		log.WithFields(log.Fields{"type": "WS_Upstream_Url"}).Trace(wssUpstreamUrl)

		dialer := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, NetDial: func(network, addr string) (net.Conn, error) {
			//强制走ipv4连接
			return net.Dial("tcp4", addr)
		}}

		upstream, _, err := dialer.Dial(wssUpstreamUrl, nil)
		if err != nil {
			log.WithFields(log.Fields{"type": "WS_Upstream_Connect"}).Warn(err.Error())
		} else {
			log.WithFields(log.Fields{"type": "WS_Upstream_Connect"}).Info(fmt.Sprintf("Upstream Websocket connect success, url:%s", wssUpstreamUrl))
		}

		go func() {
			defer downstream.Close()
			for {
				mt, payload, err := downstream.ReadMessage()
				log.WithFields(log.Fields{"type": "WS_Downstream_ReceiveMessage"}).Trace(string(payload))
				if err != nil {
					errChan <- err
					log.WithFields(log.Fields{"type": "WS_Downstream_Close"}).Info(err.Error())
					return
				}
				upload <- wsMessage{mt, payload}
			}
		}()

		go func() {
			defer upstream.Close()
			for {
				mt, payload, err := upstream.ReadMessage()
				log.WithFields(log.Fields{"type": "WS_Upstream_ReceiveMessage"}).Trace(string(payload))
				if err != nil {
					errChan <- err
					log.WithFields(log.Fields{"type": "WS_Upstream"}).Info(err.Error())
					return
				}
				download <- wsMessage{mt, payload}
			}
		}()

		go func() {
			for {
				select {
				case message := <-upload:
					err := upstream.WriteMessage(message.messageType, message.payload)
					log.WithFields(log.Fields{"type": "WS_Upstream_Forward"}).Trace(string(message.payload))
					if err != nil {
						errChan <- err
					}
				case message := <-download:
					err := downstream.WriteMessage(message.messageType, message.payload)
					log.WithFields(log.Fields{"type": "WS_Downstream_Forward"}).Trace(string(message.payload))
					if err != nil {
						errChan <- err
					}
				case err := <-errChan:
					log.WithFields(log.Fields{"type": "WS_Downstream_Close"}).Error(err.Error())
					log.WithFields(log.Fields{"type": "WS_Upstream_Close"}).Error(err.Error())
					return
				}
			}
		}()
	}
}
