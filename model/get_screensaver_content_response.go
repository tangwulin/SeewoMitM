package model

type GetScreensaverContentResponse struct {
	Response
	Data []ScreensaverContent `json:"data,omitempty"`
}
