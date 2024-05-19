package main

type SpinePlayerConfig struct {
	JsonUrl                   string            `json:"jsonUrl,omitempty"`
	SkelUrl                   string            `json:"skelUrl,omitempty"`
	AtlasUrl                  string            `json:"atlasUrl,omitempty"`
	RawDataURIs               map[string]string `json:"rawDataURIs,omitempty"`
	Animation                 string            `json:"animation,omitempty"`
	Animations                []string          `json:"animations,omitempty"`
	DefaultMix                float64           `json:"defaultMix,omitempty"`
	Skin                      string            `json:"skin,omitempty"`
	Skins                     []string          `json:"skins,omitempty"`
	PremultipliedAlpha        bool              `json:"premultipliedAlpha,omitempty"`
	ShowControls              bool              `json:"showControls,omitempty"`
	Debug                     *DebugOptions     `json:"debug,omitempty"`
	Viewport                  *ViewportOptions  `json:"viewport,omitempty"`
	Alpha                     bool              `json:"alpha,omitempty"`
	BackgroundColor           string            `json:"backgroundColor,omitempty"`
	BackgroundImage           *BackgroundImage  `json:"backgroundImage,omitempty"`
	FullScreenBackgroundColor string            `json:"fullScreenBackgroundColor,omitempty"`
	ControlBones              []string          `json:"controlBones,omitempty"`
}

type DebugOptions struct {
	Bones    bool `json:"bones,omitempty"`
	Regions  bool `json:"regions,omitempty"`
	Meshes   bool `json:"meshes,omitempty"`
	Bounds   bool `json:"bounds,omitempty"`
	Paths    bool `json:"paths,omitempty"`
	Clipping bool `json:"clipping,omitempty"`
	Points   bool `json:"points,omitempty"`
	Hulls    bool `json:"hulls,omitempty"`
}

type ViewportOptions struct {
	X              int                 `json:"x,omitempty"`
	Y              int                 `json:"y,omitempty"`
	Width          int                 `json:"width,omitempty"`
	Height         int                 `json:"height,omitempty"`
	PadLeft        int                 `json:"padLeft,omitempty"`
	PadRight       int                 `json:"padRight,omitempty"`
	PadTop         int                 `json:"padTop,omitempty"`
	PadBottom      int                 `json:"padBottom,omitempty"`
	Animations     map[string]Viewport `json:"animations,omitempty"`
	DebugRender    bool                `json:"debugRender,omitempty"`
	TransitionTime float64             `json:"transitionTime,omitempty"`
}

type Viewport struct {
	X         int `json:"x,omitempty"`
	Y         int `json:"y,omitempty"`
	Width     int `json:"width,omitempty"`
	Height    int `json:"height,omitempty"`
	PadLeft   int `json:"padLeft,omitempty"`
	PadRight  int `json:"padRight,omitempty"`
	PadTop    int `json:"padTop,omitempty"`
	PadBottom int `json:"padBottom,omitempty"`
}

type BackgroundImage struct {
	Url    string `json:"url,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height"`
}
