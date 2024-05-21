package main

import (
	"SeewoMitM/internal/log"
	"SeewoMitM/model"
	"crypto/tls"
	"encoding/json"
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

		c := &Connection{URL: r.RequestURI, UpstreamConn: upstream, DownstreamConn: downstream}

		AddConnection(c)

		go func() {
			defer downstream.Close()
			defer RemoveConnection(c)
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
			defer RemoveConnection(c)
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
	if u, exist := originalPayload["url"]; exist && u.(string) == "/displayScreenSaver" {
		log.WithFields(log.Fields{"type": "ModifyPayload"}).Info("displayScreenSaver message detected!")
		content := GetScreensaverContent()
		if len(content.ExtraPayload.ScreensaverContent) == 0 {
			return payload
		}

		var pictureSizeType int
		switch content.Fit {
		case "contain":
			pictureSizeType = 0
		case "cover":
			pictureSizeType = 1
		}

		newPayload := model.ScreensaverPayload{}
		err := json.Unmarshal(*payload, &newPayload)
		if err != nil {
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("failed to unmarshal new payload", err.Error())
			return payload
		}

		switch content.Mode {
		case "replace":
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Info("replace mode")
			// 直接替换
			newPayload.Data.ImageList = content.ImageList
			newPayload.Data.ExtraPayload = content.ExtraPayload
		case "append":
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Info("append mode")
			// 先取出原有的图片
			originalImageList := originalPayload["data"].(map[string]interface{})["imageList"].([]string)

			// 再追加新的图片
			// 对原版前端兼容
			newPayload.Data.ImageList = append(originalImageList, content.ImageList...)

			// 新版前端需要把图片转换为Content
			originalImageContent := make([]model.ScreensaverContentPayload, len(originalImageList)+len(content.ExtraPayload.ScreensaverContent))

			// 先放原来的
			for _, v := range originalImageList {
				originalImageContent = append(originalImageContent, *NewImageContent(v, content.Fit, 0))
			}

			newPayload.Data.ExtraPayload = content.ExtraPayload

			// 别忘了把其余Content也放进去
			newPayload.Data.ExtraPayload.ScreensaverContent = append(originalImageContent, content.ExtraPayload.ScreensaverContent...)
		case "off":
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Info("screensaver modify off")
		// do nothing
		default:
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("unknown mode")
		}

		newPayload.Data.PictureSizeType = pictureSizeType
		newPayload.Data.SwitchInterval = content.SwitchInterval
		newPayload.Data.MaterialSource = content.Source

		newPayloadBytes, err := json.Marshal(newPayload)
		if err != nil {
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("failed to marshal new payload", err.Error())
			return payload
		}
		return &newPayloadBytes
	}

	isDesktopAssistantMessage := false
	if messageType, exist := originalPayload["messageType"]; exist {
		switch messageType.(type) {
		case float64:
			if messageType.(float64) == 1315 {
				isDesktopAssistantMessage = true
			}
		default:
			// do nothing
		}
	}

	if isDesktopAssistantMessage {
		log.WithFields(log.Fields{"type": "ModifyPayload"}).Info("desktop assistant message detected!")
		newPayload := model.DesktopAssistantPayload{}
		err := json.Unmarshal(*payload, &newPayload)
		if err != nil {
			log.WithFields(log.Fields{"type": "ModifyPayload"}).Error("failed to unmarshal original payload", err.Error())
			return payload
		}

		newPayload.Data.Urls = []struct {
			Image  string `json:"image"`
			Name   string `json:"name"`
			Target string `json:"target"`
			Type   string `json:"type"`
		}{
			{
				Image:  "https://www.miyoushe.com/assets/ys-logo-v2-B_XJ9psI.png",
				Name:   "原神",
				Target: "https://ys.mihoyo.com/",
				Type:   "moreapps",
			},
			{
				Image:  "https://i0.hdslb.com/bfs/game/c55c98f1f31c4d2a217ffbe3187f7bce090fb6b1.jpg@280w_280h_1c_!web-search-game-cover",
				Name:   "明日方舟",
				Target: "https://ak.hypergryph.com/#index",
				Type:   "moreapps",
			},
		}

		newPayloadBytes, err := json.Marshal(newPayload)
		return &newPayloadBytes
	}

	return payload
}
