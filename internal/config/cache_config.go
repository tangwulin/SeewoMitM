package config

type CacheConfig struct {
	// 缓存大小
	Size int `json:"size" default:"1000"`

	// 缓存目录
	CacheDir string `json:"dir" default:"./cache"`

	// 缓存过期时间
	Expire int `json:"expire" default:"600"`

	// 缓存清理时间
	CleanTime int `json:"cleanTime" default:"600"`

	// 缓存清理间隔
	CleanInterval int `json:"cleanInterval" default:"600"`

	// 缓存清理百分比
	CleanPercent int `json:"cleanPercent" default:"50"`
}
