package main

import (
	"fmt"
	"reflect"
	"strings"
)

func i2s(data interface{}, out interface{}) error {
	var val reflect.Value
	if reflect.TypeOf(out).String() == "*reflect.Value" {
		ptr := out.(*reflect.Value)
		val = *ptr
	} else {
		val = reflect.ValueOf(out)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		} else {
			return fmt.Errorf("not ptr")
		}
	}
	switch val.Kind() {
	case reflect.Int:
		v, ok := data.(float64)
		if !ok {
			return fmt.Errorf("fail int")
		}
		val.Set(reflect.ValueOf(int(v)))
	case reflect.Bool:
		v, ok := data.(bool)
		if !ok {
			return fmt.Errorf("fail bool")
		}
		val.Set(reflect.ValueOf(v))
	case reflect.String:
		v, ok := data.(string)
		if !ok {
			return fmt.Errorf("fail string")
		}
		val.SetString(v)
	case reflect.Slice:
		data, ok2 := data.([]interface{})
		if !ok2 {
			return fmt.Errorf("fail slice data")
		}
		if strings.HasSuffix(val.Type().String(), "Simple") {
			val.Set(reflect.ValueOf(make([]Simple, len(data))))
		} else if strings.HasSuffix(val.Type().String(), "IDBlock") {
			val.Set(reflect.ValueOf(make([]IDBlock, len(data))))
		}
		for i := 0; i < len(data); i++ {
			v := val.Index(i)
			i2s(data[i], &v)
		}
	case reflect.Struct:
		data, ok2 := data.(map[string]interface{})
		if !ok2 {
			return fmt.Errorf("fail slice data")
		}
		for k1, v1 := range data {
			for i := 0; i < val.NumField(); i++ {
				valueField := val.Field(i)
				typeField := val.Type().Field(i)
				if k1 == typeField.Name {
					err := i2s(v1, &valueField)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}
	return nil
}
