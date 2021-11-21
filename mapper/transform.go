package mapper

import (
	"errors"
	"reflect"
	"strings"
)

type MapName struct {
	From string
	To string
}

func Map(from, to string) MapName {
	return MapName{from, to}
}

func structValue(v interface{}) (reflect.Value, error) {
	vvalue := reflect.ValueOf(v)
	if vvalue.Kind() == reflect.Ptr {
		vvalue = vvalue.Elem()
	}
	if vvalue.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("struct needed")
	}

	return vvalue, nil
}

func fieldNames(v interface{}, mapname ...MapName) (map[string]reflect.StructField, error) {
	vvalue, err := structValue(v)
	if err != nil {
		return nil, err
	}

	names := make(map[string]string)
	for _, mn := range mapname {
		names[mn.From] = mn.To
	}

	result := make(map[string]reflect.StructField)
	vtype := vvalue.Type()
	for i := 0; i < vtype.NumField(); i++ {
		fld := vtype.Field(i)
		uname := []rune(fld.Name)
		if uname[0] >= 65 && uname[0] <= 90 {
			uname[0] = uname[0] + 32
			name := string(uname)
			if tag, ok := fld.Tag.Lookup("json"); ok {
				name = strings.SplitN(tag, ",", 2)[0]
			}
			if tag, ok := fld.Tag.Lookup("mapper"); ok {
				name = strings.SplitN(tag, ",", 2)[0]
			}
			if mn, ok := names[name]; ok {
				name = mn
			}

			result[name] = fld
		}
	}

	return result, nil
}

func Transform(v interface{}, mapname ...MapName) (map[string]interface{}, error) {
	vvalue, err := structValue(v)
	if err != nil {
		return nil, err
	}

	names, err := fieldNames(v, mapname...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for name, fld := range names {
		val, err := transformValue(vvalue.FieldByName(fld.Name))
		if err != nil {
			return nil, err
		}
		result[name] = val
	}

	return result, nil
}

func transformValue(value reflect.Value) (interface{}, error) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		return transformSliceValue(value)
	case reflect.Map:
		return transformMapValue(value)
	case reflect.Struct:
		return Transform(value)
	case reflect.Chan:
		return nil, errors.New("chan not supported")
	}
	return value.Interface(), nil
}

func transformMapValue(value reflect.Value) (map[string]interface{}, error) {
	if value.Type().Key().Kind() != reflect.String {
		return nil, errors.New("only maps with string keys supported")		
	}

	result := make(map[string]interface{})
	for _, key := range value.MapKeys() {
		val, err := transformValue(value.MapIndex(key))
		if err != nil {
			return nil, err
		}
		result[key.String()] = val
	}

	return result, nil
}

func transformSliceValue(value reflect.Value) ([]interface{}, error) {
	result := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		val, err := transformValue(value.Index(i))
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}


// TODO Convert func to fill struct-pointer