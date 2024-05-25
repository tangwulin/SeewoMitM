package model

type AddScreensaverContentResponse struct {
	Response
	Data int `json:"data,omitempty"` //添加后得到的id
}
