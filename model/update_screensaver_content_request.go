package model

type UpdateScreensaverContentRequest struct {
	ID int `json:"id"`
	// The screensaver content to be updated.
	Content *ScreensaverContent `json:"content"`
}
