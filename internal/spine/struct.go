package spine

type Options struct {
	JsonUrl                   string            `json:"jsonUrl"`
	SkelUrl                   string            `json:"skelUrl"`
	AtlasUrl                  string            `json:"atlasUrl"`
	RawDataURIs               map[string]string `json:"rawDataURIs"`
	Animation                 string            `json:"animation"`
	Animations                []string          `json:"animations"`
	DefaultMix                float64           `json:"defaultMix"`
	Skin                      string            `json:"skin"`
	Skins                     []string          `json:"skins"`
	PremultipliedAlpha        bool              `json:"premultipliedAlpha"`
	ShowControls              bool              `json:"showControls"`
	Debug                     DebugOptions      `json:"debug"`
	Viewport                  ViewportOptions   `json:"viewport"`
	Alpha                     bool              `json:"alpha"`
	BackgroundColor           string            `json:"backgroundColor"`
	BackgroundImage           BackgroundImage   `json:"backgroundImage"`
	FullScreenBackgroundColor string            `json:"fullScreenBackgroundColor"`
	ControlBones              []string          `json:"controlBones"`
}

type DebugOptions struct {
	Bones    bool `json:"bones"`
	Regions  bool `json:"regions"`
	Meshes   bool `json:"meshes"`
	Bounds   bool `json:"bounds"`
	Paths    bool `json:"paths"`
	Clipping bool `json:"clipping"`
	Points   bool `json:"points"`
	Hulls    bool `json:"hulls"`
}

type ViewportOptions struct {
	X              int                 `json:"x"`
	Y              int                 `json:"y"`
	Width          int                 `json:"width"`
	Height         int                 `json:"height"`
	PadLeft        int                 `json:"padLeft"`
	PadRight       int                 `json:"padRight"`
	PadTop         int                 `json:"padTop"`
	PadBottom      int                 `json:"padBottom"`
	Animations     map[string]Viewport `json:"animations"`
	DebugRender    bool                `json:"debugRender"`
	TransitionTime float64             `json:"transitionTime"`
}

type Viewport struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Width     int `json:"width"`
	Height    int `json:"height"`
	PadLeft   int `json:"padLeft"`
	PadRight  int `json:"padRight"`
	PadTop    int `json:"padTop"`
	PadBottom int `json:"padBottom"`
}

type BackgroundImage struct {
	Url    string `json:"url"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
