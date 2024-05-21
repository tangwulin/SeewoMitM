package model

type ScreensaverPayload struct {
	Data    ScreensaverPayloadData `json:"data"`
	TraceId string                 `json:"traceId"`
	Url     string                 `json:"url"`
}

type ScreensaverPayloadData struct {
	ImageList       []string      `json:"imageList"`
	MaterialSource  string        `json:"materialSource"`
	ExtraPayload    *ExtraPayload `json:"extraPayload,omitempty"`
	PictureSizeType int           `json:"pictureSizeType"`
	PlayMode        int           `json:"playMode"`
	SwitchInterval  int           `json:"switchInterval"`
	TextList        []TextItem    `json:"textList"`
}

type ExtraPayload struct {
	ScreensaverSwitchInterval int                         `json:"screensaverSwitchInterval"`
	ScreensaverContent        []ScreensaverContentPayload `json:"screensaverContent"`
}

type TextItem struct {
	Content    string `json:"content"`
	Provenance string `json:"provenance"`
}
