package config

type MitMConfig struct {
	Enable bool `json:"enable"`
	Rules  []struct {
		URL         string `json:"url"`
		Direction   string `json:"direction"`
		MessageType int    `json:"messageType"`
		Action      []struct {
			Type  string      `json:"type"`
			Key   string      `json:"key,omitempty"`
			Value interface{} `json:"value,omitempty"`
		}
	} `json:"rules"`
}
