package model

type GetScreensaverContentByIDResponse struct {
	Response
	Data ScreensaverContent `json:"data,omitempty"`
}
