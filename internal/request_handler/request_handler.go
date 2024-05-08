package request_handler

import (
	"SeewoMitM/internal/connection"
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/screensaver"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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
		closeChan := make(chan error, 20)

		downstream, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.WithFields(log.Fields{"type": "WS_Downstream_Upgrade"}).Errorf("Downstream Websocket upgrade failed, url:%s", r.RequestURI)
			return
		} else {
			log.WithFields(log.Fields{"type": "WS_Downstream_Upgrade"}).Tracef("Downstream Websocket upgrade success, url:%s", r.RequestURI)
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
			downstream.Close()
			return
		} else {
			log.WithFields(log.Fields{"type": "WS_Upstream_Connect"}).Tracef("Upstream Websocket connect success, url:%s", wssUpstreamUrl)
		}

		c := &connection.Connection{URL: r.RequestURI, UpstreamConn: upstream, DownstreamConn: downstream}

		connection.AddConnection(c)

		go func() {
			defer downstream.Close()
			defer connection.RemoveConnection(c)
			for {
				mt, payload, err := downstream.ReadMessage()
				log.WithFields(log.Fields{"type": "WS_Downstream_ReceiveMessage"}).Trace(string(payload))
				if err != nil {
					closeChan <- err
					log.WithFields(log.Fields{"type": "WS_Downstream_Close"}).Info(err.Error())
					return
				}
				upload <- wsMessage{mt, payload}
			}
		}()

		go func() {
			defer upstream.Close()
			defer connection.RemoveConnection(c)
			for {
				mt, payload, err := upstream.ReadMessage()
				log.WithFields(log.Fields{"type": "WS_Upstream_ReceiveMessage"}).Trace(string(payload))
				if err != nil {
					closeChan <- err
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
						closeChan <- err
						upstream.Close()
					}
				case message := <-download:
					var newPayload *[]byte
					if r.RequestURI == "/forward/SeewoHugoHttp/SeewoHugoService" {
						newPayload = ModifyPayload(&message.payload)
						log.WithFields(log.Fields{"type": "WS_Downstream_Forward"}).Warn(string(*newPayload))
					} else {
						newPayload = &message.payload
					}
					err := downstream.WriteMessage(message.messageType, *newPayload)
					log.WithFields(log.Fields{"type": "WS_Downstream_Forward"}).Trace(string(*newPayload))
					if err != nil {
						closeChan <- err
						downstream.Close()
					}
				case _ = <-closeChan:
					return
				}
			}
		}()
	}
}

func ModifyPayload(payload *[]byte) *[]byte {
	originalPayload := make(map[string]interface{})
	err := json.Unmarshal(*payload, &originalPayload)
	if err != nil {
		log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("failed to unmarshal payload")
		return payload
	}

	// 屏保
	if url, exist := originalPayload["url"]; exist && url == "/displayScreenSaver" {
		content := screensaver.GetScreensaverContent()
		var pictureSizeType int
		switch content.Fit {
		case "contain":
			pictureSizeType = 0
		case "cover":
			pictureSizeType = 1
		}

		if len(content.ExtraPayload.ScreensaverContent) == 0 {
			return payload
		}

		newPayload := screensaver.Payload{}
		err := json.Unmarshal(*payload, &newPayload)
		if err != nil {
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("failed to unmarshal new payload", err.Error())
			return payload
		}

		switch content.Mode {
		case "replace":
			// 直接替换
			newPayload.Data.ImageList = content.ImageList
			newPayload.Data.ExtraPayload = content.ExtraPayload
		case "append":
			// 先取出原有的图片
			originalImageList := originalPayload["data"].(map[string]interface{})["imageList"].([]string)

			// 再追加新的图片
			// 对原版前端兼容
			newPayload.Data.ImageList = append(originalImageList, content.ImageList...)

			// 新版前端用
			var originalImageContent []screensaver.Content
			for _, v := range originalImageList {
				originalImageContent = append(originalImageContent, *screensaver.NewImageContent(v, content.Fit, 0))
			}
			newPayload.Data.ExtraPayload = content.ExtraPayload
		case "off":
		// do nothing
		default:
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("unknown mode")
		}

		newPayload.Data.PictureSizeType = pictureSizeType
		newPayload.Data.SwitchInterval = content.SwitchInterval
		newPayload.Data.MaterialSource = content.Source
		newPayload.TraceId = strings.ToUpper(uuid.New().String())

		newPayloadBytes, err := json.Marshal(newPayload)
		if err != nil {
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("failed to marshal new payload", err.Error())
			return payload
		}
		return &newPayloadBytes
	}

	if messageType, exist := originalPayload["messageType"]; exist && messageType == 1315 {
		// TODO:
	}

	return payload
}
