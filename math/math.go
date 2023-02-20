// Package math 计算实用工具
package math

import "math"

type num interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// Max 返回两数最大值, 该函数将被内联
func Max[T num](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min 返回两数最小值, 该函数将被内联
func Min[T num](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// Ceil 向上整除(除数,被除数)
func Ceil[T num](dividend, divisor T) T {
	result := dividend / divisor
	if dividend%divisor != 0 {
		return result + 1
	}
	return result
}

// intSize is either 32 or 64.
const intSize = 32 << (^uint(0) >> 63)

// Abs 返回绝对值, 该函数将被内联
func Abs(x int) int {
	// m := -1 if x < 0. m := 0 otherwise.
	m := x >> (intSize - 1)

	// In two's complement representation, the negative number
	// of any number (except the smallest one) can be computed
	// by flipping all the bits and add 1. This is faster than
	// code with a branch.
	// See Hacker's Delight, section 2-4.
	return (x ^ m) - m
}

// Abs64 返回绝对值, 该函数将被内联
func Abs64(x int64) int64 {
	// m := -1 if x < 0. m := 0 otherwise.
	m := x >> (64 - 1)

	// In two's complement representation, the negative number
	// of any number (except the smallest one) can be computed
	// by flipping all the bits and add 1. This is faster than
	// code with a branch.
	// See Hacker's Delight, section 2-4.
	return (x ^ m) - m
}

// Partition 分片函数
func Partition[T any](list []T, size int) (plist [][]T) {
	if size <= 0 {
		return nil
	}
	plen := int(math.Ceil(float64(len(list)) / float64(size)))
	plist = make([][]T, plen)
	for i := 0; i < plen-1; i++ {
		plist[i] = list[i*size : (i+1)*size]
	}
	plist[plen-1] = list[(plen-1)*size:]
	return
}
