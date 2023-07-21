package utils

import (
	"fmt"
	"reflect"
)

// 获取类型信息：reflect.TypeOf，是静态的
// 获取值信息：reflect.ValueOf，是动态的

// KeyBy 提取Slice结构体的某个字段转换成Map，key field , val Struct
func KeyBy(list interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 校验传入数据的类型
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		panic(any("list required slice or array type"))
	}
	// 拿到数组的类型
	valType := val.Type()

	// 拿到数组的元素的类型
	elemType := valType.Elem()

	// elemType 用于声明 map 的 value 类型
	// elemValType 用于拿到结构体里对应字段的类型
	elemValType := elemType

	// 指针特殊处理
	for elemValType.Kind() == reflect.Ptr {
		elemValType = elemValType.Elem()
	}

	// 校验是否结构体
	if elemValType.Kind() != reflect.Struct {
		panic(any("element not struct"))
	}

	// 获取字段
	field, ok := elemValType.FieldByName(fieldName)
	if !ok {
		panic(any(fmt.Sprintf("field %s not found", fieldName)))
	}

	// 初始化存储的map
	resultMap := reflect.MakeMapWithSize(reflect.MapOf(field.Type, elemType), val.Len())

	// range slice or array set value to map
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		elemStruct := elem
		for elemStruct.Kind() == reflect.Ptr {
			elemStruct = elemStruct.Elem()
		}

		// 如果是nil的，意味着key和value同时不存在，所以跳过不处理
		if !elemStruct.IsValid() {
			continue
		}

		if elemStruct.Kind() != reflect.Struct {
			panic(any("element not struct"))
		}

		resultMap.SetMapIndex(elemStruct.FieldByIndex(field.Index), elem)
	}

	return resultMap.Interface()
}
