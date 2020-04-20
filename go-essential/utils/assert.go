package utils

import (
	"strconv"
)

func GetAssertString(data interface{}) string {
	if str, ok := data.(string); ok {
		return str
	} else if intValue, ok := data.(int); ok {
		return strconv.Itoa(intValue)
	} else if intValue, ok := data.(int32); ok {
		return strconv.FormatInt(int64(intValue), 10)
	} else if intValue, ok := data.(float64); ok {
		return strconv.FormatFloat(intValue, 'f', -1, 64)
	} else if intValue, ok := data.(int64); ok {
		return strconv.FormatInt(intValue, 10)
	}

	return ""
}

func GetAssertInt(data interface{}) int {
	var value int
	switch data.(type) {
	case string:
		value, _ = strconv.Atoi(data.(string))
		break
	case int:
		value = data.(int)
		break
	case int32:
		value = int(data.(int32))
		break
	case int64:
		value = int(data.(int64))
		break
	case float64:
		value = int(data.(float64))
		break
	case float32:
		value = int(data.(float32))
		break
	default:
		value = 0
		break
	}

	return value
}

func GetAssertInt32(data interface{}) int32 {
	var value int32
	switch data.(type) {
	case string:
		if val, err := strconv.Atoi(data.(string)); err == nil {
			value = int32(val)
		} else {
			value = 0
		}
		break
	case int:
		value = int32(data.(int))
		break
	case int32:
		value = data.(int32)
		break
	case int64:
		value = int32(data.(int64))
		break
	case float64:
		value = int32(data.(float64))
		break
	case float32:
		value = int32(data.(float32))
		break
	default:
		//log.Debug("data type:", reflect.TypeOf(data))
		value = 0
		break
	}

	return value
}

func GetAssertInt64(data interface{}) int64 {
	var value int64
	switch data.(type) {
	case string:
		if val, err := strconv.Atoi(data.(string)); err == nil {
			value = int64(val)
		} else {
			value = 0
		}
		break
	case int:
		value = int64(data.(int))
		break
	case int32:
		value = int64(data.(int32))
		break
	case float64:
		value = int64(data.(float64))
		break
	case float32:
		value = int64(data.(float32))
		break
	case int64:
		value = data.(int64)
		break
	default:
		value = 0
		break
	}

	return value
}

func GetAssertBool(data interface{}) bool {
	if data == nil {
		return false
	}
	switch data.(type) {
	case bool:
		return data.(bool)
	}

	return false
}

func GetAssertFloat64(data interface{}) float64 {
	var f float64
	if str, ok := data.(string); ok {
		float, err := strconv.ParseFloat(str, 64)
		if err == nil {
			return float
		}
	} else if intValue, ok := data.(int); ok {
		return float64(intValue)
	} else if intValue, ok := data.(int32); ok {
		return float64(intValue)
	} else if intValue, ok := data.(float64); ok {
		return intValue
	} else if intValue, ok := data.(int64); ok {
		return float64(intValue)
	}

	return f
}
