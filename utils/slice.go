package utils

import (
	"fmt"
	"math/rand"
	"reflect"
)

// InterfaceSlice 将普通数组转换为interface数组
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Uint32InterfaceSlice 将interface数组转换为uint32数组
func Uint32InterfaceSlice(slice []interface{}) []uint32 {
	ret := make([]uint32, len(slice))
	for i, v := range slice {
		ret[i] = v.(uint32)
	}
	return ret
}

// IntInterfaceSlice 将interface数组转换为int32数组
func IntInterfaceSlice(slice []interface{}) []int {
	ret := make([]int, len(slice))
	for i, v := range slice {
		ret[i] = v.(int)
	}
	return ret
}

// StringInterfaceSlice 将interface数组转换为string数组
func StringInterfaceSlice(slice []interface{}) []string {
	ret := make([]string, len(slice))
	for i, v := range slice {
		ret[i] = v.(string)
	}
	return ret
}

// SliceContains 元素是否在数组里面(注意类型一定要一样)
func SliceContains(slice interface{}, element interface{}) bool {
	for _, v := range InterfaceSlice(slice) {
		if v == element {
			return true
		}
	}
	return false
}

// SliceContainsUint16 元素是否在数组里面Uint16
func SliceContainsUint16(slice []uint16, element uint16) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// SliceContainsUint32 元素是否在数组里面Uint32
func SliceContainsUint32(slice []uint32, element uint32) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// SliceContainsUint64 元素是否在数组里面
func SliceContainsUint64(slice []uint64, element uint64) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// SliceContainsInt 元素是否在数组里面
func SliceContainsInt(slice []int, element int) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// SliceContainsString 元素是否在数组里面
func SliceContainsString(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// SliceRemoveDuplicate
/* 在slice中去除重复的元素，其中a必须是已经排序的序列。
 * params:
 *   a: slice对象，如[]string, []int, []float64, ...
 * return:
 *   []interface{}: 已经去除重复元素的新的slice对象
 */
func SliceRemoveDuplicate(a interface{}) (ret []interface{}) {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		fmt.Printf("<SliceRemoveDuplicate> <a> is not slice but %T\n", a)
		return ret
	}

	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}

	return ret
}

// SliceInsert 往数组插入元素
func SliceInsert(s []interface{}, index int, value interface{}) []interface{} {
	rear := append([]interface{}{}, s[index:]...)
	return append(append(s[:index], value), rear...)
}

// SliceDifference 差集：以属于A而不属于B的元素为元素的集合成为A与B的差 (集)
func SliceDifference(slice1, slice2 interface{}) []interface{} {
	var slice3 []interface{}
	for _, v := range InterfaceSlice(slice1) {
		if !SliceContains(slice2, v) {
			slice3 = append(slice3, v)
		}
	}
	return slice3
}

// SliceIntersect 交集： 以属于A且属于B的元素为元素的集合成为A与B的交（集）
func SliceIntersect(slice1, slice2 interface{}) []interface{} {
	var slice3 []interface{}
	for _, v := range InterfaceSlice(slice1) {
		if SliceContains(slice2, v) {
			slice3 = append(slice3, v)
		}
	}
	return slice3
}

// SliceUnion 并集：以属于A或属于B的元素为元素的集合成为A与B的并（集）
func SliceUnion(slice1, slice2 interface{}) []interface{} {
	var slice3 []interface{}
	for _, v := range InterfaceSlice(slice1) {
		if SliceContains(slice2, v) {
			slice3 = append(slice3, v)
		}
	}
	slice3 = append(slice3, slice2.([]interface{})...)
	return slice3
}

// SliceFind 传入方法找到元素
func SliceFind(slice interface{}, findFunc func(element interface{}, idx int) bool) interface{} {
	for i, v := range InterfaceSlice(slice) {
		if findFunc(v, i) {
			return v
		}
	}

	return nil
}

// SliceFindIndex 传入方法找到元素索引
func SliceFindIndex(slice interface{}, findFunc func(element interface{}, idx int) bool) int {
	for i, v := range InterfaceSlice(slice) {
		if findFunc(v, i) {
			return i
		}
	}

	return -1
}

// RandSlice 切片乱序
func RandSlice(slice interface{}) {
	rv := reflect.ValueOf(slice)
	if rv.Type().Kind() != reflect.Slice {
		return
	}

	length := rv.Len()
	if length < 2 {
		return
	}

	swap := reflect.Swapper(slice)
	for i := length - 1; i >= 0; i-- {
		j := rand.Intn(length)
		swap(i, j)
	}
	return
}
