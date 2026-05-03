package utils

import (
	"reflect"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 跳过未导出字段
		if fieldType.PkgPath != "" {
			continue
		}

		fieldName := fieldType.Name

		// 递归处理嵌套结构体
		if field.Kind() == reflect.Struct {
			result[fieldName] = Struct2Map(field.Interface())
		} else if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if !field.IsNil() {
				result[fieldName] = Struct2Map(field.Elem().Interface())
			} else {
				result[fieldName] = nil
			}
		} else {
			result[fieldName] = field.Interface()
		}
	}

	return result
}
