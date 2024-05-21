package model

type DesktopAssistantPayload struct {
	Data struct {
		Apps []struct {
			ExePath string `json:"exePath"`
			Image   string `json:"image"`
			Name    string `json:"name"`
			Type    string `json:"type"`
		} `json:"apps"`
		Urls []struct {
			Image  string `json:"image"`
			Name   string `json:"name"`
			Target string `json:"target"`
			Type   string `json:"type"`
		} `json:"urls"`
		Visible bool `json:"visible"`
	} `json:"data"`
	MessageType int    `json:"messageType"`
	TraceId     string `json:"traceId"`
}
