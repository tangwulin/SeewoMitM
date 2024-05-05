package screensaver

import "github.com/google/uuid"

type Payload struct {
	Data    Data   `json:"data"`
	TraceId string `json:"traceId"`
	Url     string `json:"url"`
}

type Data struct {
	ImageList       []string      `json:"imageList"`
	MaterialSource  string        `json:"materialSource"`
	ExtraPayload    *ExtraPayload `json:"extraPayload,omitempty"`
	PictureSizeType int           `json:"pictureSizeType"`
	PlayMode        int           `json:"playMode"`
	SwitchInterval  int           `json:"switchInterval"`
	TextList        []TextItem    `json:"textList"`
}

type ExtraPayload struct {
	ScreensaverContent []Content `json:"screensaverContent"`
}

type TextItem struct {
	Content    string `json:"content"`
	Provenance string `json:"provenance"`
}

func NewPayload(imageList []string, materialSource string, extraPayload *ExtraPayload, pictureSizeType int, playMode int, switchInterval int, textList []TextItem) *Payload {
	return &Payload{
		Data: Data{
			ImageList:       imageList,
			MaterialSource:  materialSource,
			ExtraPayload:    extraPayload,
			PictureSizeType: pictureSizeType,
			PlayMode:        playMode,
			SwitchInterval:  switchInterval,
			TextList:        textList,
		},
		TraceId: uuid.New().String(),
		Url:     "/displayScreenSaver",
	}
}
