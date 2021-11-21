package mapper

import (
	"strings"
	"time"
)

func Value(path string, v interface{}, mapname ...MapName) (interface{}, bool) {
	// TODO also support maps
	if len(path) == 0 {
		return nil, false
	}

	vvalue, err := structValue(v)
	if err != nil {
		return nil, false
	}

	names, err := fieldNames(v, mapname...)
	if err != nil {
		return nil, false
	}

	split := strings.SplitN(path, ".", 2)
	name := split[0] // TODO slice[index]
	pth := ""
	if len(split) > 1 {
		pth = split[1]
	}

	if fld, ok := names[name]; ok {
		val := vvalue.FieldByName(fld.Name).Interface()
		if len(pth) > 0 {
			return Value(pth, val, mapname...)
		}
		return val, true
	}
	
	return nil, false
}

func IntValue(path string, v interface{}, mapname ...MapName) (int, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(int); ok {
			return i, true
		}
	}
	return 0, false
}

func Int8Value(path string, v interface{}, mapname ...MapName) (int8, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(int8); ok {
			return i, true
		}
	}
	return 0, false
}

func Int16Value(path string, v interface{}, mapname ...MapName) (int16, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(int16); ok {
			return i, true
		}
	}
	return 0, false
}

func Int32Value(path string, v interface{}, mapname ...MapName) (int32, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(int32); ok {
			return i, true
		}
	}
	return 0, false
}

func Int64Value(path string, v interface{}, mapname ...MapName) (int64, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		switch ix := v.(type) {
		case int:
			return int64(ix), true
		case int8:
			return int64(ix), true
		case int16:
			return int64(ix), true
		case int32:
			return int64(ix), true
		case int64:
			return ix, true
		case uint:
			return int64(ix), true
		case uint8:
			return int64(ix), true
		case uint16:
			return int64(ix), true
		case uint32:
			return int64(ix), true
		case uint64:
			return int64(ix), true
		}
	}
	return 0, false
}

func UintValue(path string, v interface{}, mapname ...MapName) (uint, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(uint); ok {
			return i, true
		}
	}
	return 0, false
}

func Uint8Value(path string, v interface{}, mapname ...MapName) (uint8, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(uint8); ok {
			return i, true
		}
	}
	return 0, false
}

func Uint16Value(path string, v interface{}, mapname ...MapName) (uint16, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(uint16); ok {
			return i, true
		}
	}
	return 0, false
}

func Uint32Value(path string, v interface{}, mapname ...MapName) (uint32, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(uint32); ok {
			return i, true
		}
	}
	return 0, false
}

func Uint64Value(path string, v interface{}, mapname ...MapName) (uint64, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(uint64); ok {
			return i, true
		}
	}
	return 0, false
}

func Float32Value(path string, v interface{}, mapname ...MapName) (float32, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(float32); ok {
			return i, true
		}
	}
	return 0, false
}

func Float64Value(path string, v interface{}, mapname ...MapName) (float64, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		switch fx := v.(type) {
		case float32:
			return float64(fx), true
		case float64:
			return fx, true
		}
	}
	return 0, false
}

func BoolValue(path string, v interface{}, mapname ...MapName) (bool, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(bool); ok {
			return i, true
		}
	}
	return false, false
}

func StringValue(path string, v interface{}, mapname ...MapName) (string, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if i, ok := v.(string); ok {
			return i, true
		}
	}
	return "", false
}

func TimeValue(path string, v interface{}, mapname ...MapName) (time.Time, bool) {
	if v, ok := Value(path, v, mapname...); ok {
		if t, ok := v.(time.Time); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

// TODO SetXXXValue 