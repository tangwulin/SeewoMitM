package model

type DataContent struct {
	Mode           string        `json:"mode"`
	ImageList      []string      `json:"imageList"`
	ExtraPayload   *ExtraPayload `json:"extraPayload,omitempty"`
	Source         string        `json:"source"`
	Fit            string        `json:"fit"`
	PlayMode       string        `json:"playMode"`
	SwitchInterval int           `json:"switchInterval"`
	TextList       []TextItem    `json:"textList"`
}
