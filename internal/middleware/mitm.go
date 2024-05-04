package middleware

import (
	"SeewoMitM/internal/config"
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/screensaver"
	"encoding/json"
)

func ScreenSaverMitM(url string, direction string, messageType int, payload []byte) []byte {
	if url != "/forward/SeewoHugoHttp/SeewoHugoService" {
		return payload
	}
	data := make(map[string]interface{})
	if urlInData, exist := data["url"]; exist && urlInData == "/displayScreenSaver" {
		c := helper.GetConfig()
		if c == nil {
			return payload
		}

		type extraPayload struct {
			ScreenSaverContent []screensaver.ScreenSaverContent `json:"screenSaverContent"`
		}

		switch c.ScreenSaverHijackMode {
		case config.ScreenSaverHijackModeReplaceAll:
		case config.ScreenSaverHijackModeAdd:
			o, _ := data["imageList"]
			s := helper.InterfaceToSlice(o)

			// 在这里咒骂一下golang的类型系统
			originalList := make([]string, 0)

			if c.ScreenSaverHijackMode == config.ScreenSaverHijackModeAdd {
				for _, val := range s {
					originalList = append(originalList, val.(string))
				}
			}

			for _, val := range screensaver.GetScreenSaverContent() {
				if val.Type == "image" {
					originalList = append(originalList, val.Url)
				}
			}

			data["extraPayload"] = screensaver.GetScreenSaverContent()

			newPayload, err := json.Marshal(data)
			if err != nil {
				return payload
			}

			return newPayload
		default:
			return payload
		}
	}
	return payload
}

func WebsocketMitM(url string, direction string, messageType int, payload []byte) []byte {
	c := helper.GetConfig()
	if c == nil {
		return payload
	}

	if c.MitM == nil {
		return payload
	}

	if !c.MitM.Enable {
		return payload
	}

	for _, v := range c.MitM.Rules {
		if v.URL == url && v.Direction == direction && v.MessageType == messageType {
			for _, action := range v.Action {
				switch action.Type {
				case "replace":
				case "add":
					data := make(map[string]interface{})
					err := json.Unmarshal(payload, &data)
					if err != nil {
						return payload
					}
					data[action.Key] = action.Value
					newPayload, err := json.Marshal(data)
					if err != nil {
						return payload
					}
					return newPayload
				case "delete":
					data := make(map[string]interface{})
					err := json.Unmarshal(payload, &data)
					if err != nil {
						return payload
					}
					data[action.Key] = nil
					newPayload, err := json.Marshal(data)
					if err != nil {
						return payload
					}
					return newPayload
				case "reject":
					return nil
				}
			}
		}
	}

	return payload
}
