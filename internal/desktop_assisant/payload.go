package desktop_assisant

type Payload struct {
	Data struct {
		Apps []interface{} `json:"apps"`
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
