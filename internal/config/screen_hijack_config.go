package config

import "SeewoMitM/internal/spine"

type ScreenSaverHijackMode int32

type ScreenSaverHijackContent struct {
	Type           ScreenSaverHijackContentType `json:"type"`
	Path           string                       `json:"path,omitempty"`
	RequirePreload bool                         `json:"requirePreload,omitempty"`
	EntryPoint     string                       `json:"entryPoint,omitempty"`
	ServePath      string                       `json:"servePath,omitempty"`
	SpineVersion   string                       `json:"spineVersion,omitempty"`
	SpineConfig    *spine.Options               `json:"spineConfig,omitempty"`
}

const (
	ScreenSaverHijackModeOff ScreenSaverHijackMode = iota
	ScreenSaverHijackModeAdd
	ScreenSaverHijackModeReplaceAll
)

type ScreenSaverHijackContentType string

const (
	HTMLScreenSaverHijackContent           ScreenSaverHijackContentType = "html"
	ImageScreenSaverHijackContent          ScreenSaverHijackContentType = "image"
	VideoScreenSaverHijackContent          ScreenSaverHijackContentType = "video"
	SpineScreenSaverHijackContent          ScreenSaverHijackContentType = "spine"
	ImageDirectoryScreenSaverHijackContent ScreenSaverHijackContentType = "imageDirectory"
)
