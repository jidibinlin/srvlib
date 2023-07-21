package utils

import (
	"fmt"
	"github.com/gzjjyz/srvlib/utils/pie"
	"reflect"
)

func pluckFieldList(list interface{}, fieldName string) (result []reflect.Value) {
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			// 拿到元素
			elem := val.Index(i)
			// 指针需要进一步转换
			if elem.Kind() == reflect.Ptr {
				// Elem()获取地址指向的值
				elem = elem.Elem()
			}
			// 判断是否为空，为空就跳过
			if !elem.IsValid() {
				continue
			}
			// 判断是否为结构体
			if elem.Kind() != reflect.Struct {
				panic(any("element not struct"))
			}
			// 通过字段名称拿到这个值
			f := elem.FieldByName(fieldName)
			if !f.IsValid() {
				panic(any(fmt.Sprintf("struct missed field %s", fieldName)))
			}
			result = append(result, f)
		}
	default:
		panic(any("required list of struct type"))
	}
	return
}

func PluckStrings(list interface{}, fieldName string) pie.Strings {
	var result []string
	l := pluckFieldList(list, fieldName)
	for _, f := range l {
		// 判断值的类型
		if f.Kind() != reflect.String {
			panic(any(fmt.Sprintf("struct element %s type required int", fieldName)))
		}
		// 加入list中
		result = append(result, f.String())
	}
	return result
}

func PluckUint32s(list interface{}, fieldName string) pie.Uint32s {
	var result []uint32
	l := pluckFieldList(list, fieldName)
	for _, f := range l {
		// 判断值的类型
		if f.Kind() != reflect.Uint32 {
			panic(any(fmt.Sprintf("struct element %s type required int", fieldName)))
		}
		// 加入list中
		result = append(result, uint32(f.Uint()))
	}
	return result
}
func PluckUint64s(list interface{}, fieldName string) pie.Uint64s {
	var result []uint64
	l := pluckFieldList(list, fieldName)
	for _, f := range l {
		// 判断值的类型
		if f.Kind() != reflect.Uint64 {
			panic(any(fmt.Sprintf("struct element %s type required int", fieldName)))
		}
		// 加入list中
		result = append(result, f.Uint())
	}
	return result
}

func PluckInt32s(list interface{}, fieldName string) pie.Int32s {
	var result []int32
	l := pluckFieldList(list, fieldName)
	for _, f := range l {
		// 判断值的类型
		if f.Kind() != reflect.Int32 {
			panic(any(fmt.Sprintf("struct element %s type required int", fieldName)))
		}
		// 加入list中
		result = append(result, int32(f.Int()))
	}
	return result
}
func PluckInt64s(list interface{}, fieldName string) pie.Int64s {
	var result []int64
	l := pluckFieldList(list, fieldName)
	for _, f := range l {
		// 判断值的类型
		if f.Kind() != reflect.Int64 {
			panic(any(fmt.Sprintf("struct element %s type required int", fieldName)))
		}
		// 加入list中
		result = append(result, f.Int())
	}
	return result
}
