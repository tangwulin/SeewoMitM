package helper

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"net/http"
)

func TestWSAvailable(urlStr string, requestHeader http.Header) bool {
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	ws, _, err := dialer.Dial(urlStr, requestHeader)
	if err != nil {
		return false
	}
	defer ws.Close()

	return true
}
