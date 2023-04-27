package random

import (
	"github.com/gzjjyz/srvlib/utils"
	"math/rand"
)

// UintU 生成一个1-n的随机数
func UintU(n uint32) uint32 {
	if n == 0 {
		return 0
	}
	return uint32(rand.Intn(int(n)) + 1)
}

// Interval 生成指定[]区间一个随机数
func Interval(low, high int) int {
	if low == high {
		return low
	}
	if low > high {
		low, high = high, low
	}
	return rand.Intn(high-low+1) + low
}

// IntervalU 生成指定区间一个随机数
func IntervalU(low, high int) uint32 {
	return uint32(Interval(low, high))
}

func IntervalUU(low, high uint32) uint32 {
	return IntervalU(int(low), int(high))
}

// Interval64 生成指定[]区间一个随机数 不能为负数
func Interval64(low, high int64) int64 {
	if low < 0 {
		low = 0
	}
	if high < 0 {
		high = 0
	}
	if low == high {
		return low
	}
	if low > high {
		low, high = high, low
	}
	return rand.Int63n(high-low+1) + low
}

// RandPerm 生成一个随机序列索引
func RandPerm(n int) []int {
	return rand.Perm(n)
}

// RandMany 区间内随机多个不重复值
func RandMany(min, max, count uint32) []uint32 {
	if min > max {
		min, max = max, min
	}
	if count == 0 {
		count = 1
	}

	tmp := make(map[uint32]uint32)
	size := min + count - 1
	size = utils.MaxUInt32(size, max-min-1)
	ret := make([]uint32, count)
	idx := 0
	for i := min; i <= size; i++ {
		j := IntervalUU(i, max)
		if _, ok := tmp[i]; !ok {
			tmp[i] = i
		}
		if _, ok := tmp[j]; !ok {
			tmp[j] = j
		}
		tmp[i], tmp[j] = tmp[j], tmp[i]
		ret[idx] = tmp[i]
		idx++
		if idx >= int(count) {
			return ret
		}
	}
	return ret
}

// RandString 生成指定长度的随机字符串
func RandString(n int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

const (
	keyBase = 48
	keyRand = 122 - 48 + 1
)

func GenerateKey(n int) string {
	bytes := make([]byte, n)
	// key:ascII 48- 122
	for i := 0; i < n; i++ {
		bytes[i] = byte(rand.Intn(keyRand) + keyBase)
	}
	return string(bytes)
}
