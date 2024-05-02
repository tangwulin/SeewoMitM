package mitm

import (
	"SeewoMitM/internal/helper"
	"encoding/json"
)

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
