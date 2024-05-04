package helper

func MapToSlice(m map[string]interface{}) []interface{} {
	s := make([]interface{}, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func InterfaceToSlice(i interface{}) []interface{} {
	switch i.(type) {
	case map[string]interface{}:
		return MapToSlice(i.(map[string]interface{}))
	case []interface{}:
		return i.([]interface{})
	default:
		return []interface{}{}
	}
}
