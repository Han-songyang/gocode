package main

import "errors"

func Delete[T any](arr []T, index int) ([]T, error) {
	length := len(arr)
	if index < 0 || index >= length {
		return nil, errors.New("index out of range")
	}
	// move
	for i := index; i+1 < length; i++ {
		arr[i] = arr[i+1]
	}
	// delete last element
	arr = arr[:length-1]
	// TODO shrink
	arr = Shrink(arr)
	return arr, nil
}

func Shrink[T any](arr []T) []T {
	cap := cap(arr)
	len := len(arr)
	if len < cap/2 {
		out := make([]T, 0, len)
		out = append(out, arr...)
		return out
	}
	return arr
}

// answer
//// Shrink 这是缩容
//func Shrink[T any](src []T) []T {
//	c, l := cap(src), len(src)
//	n, changed := calCapacity(c, l)
//	if !changed {
//		return src
//	}
//	s := make([]T, 0, n)
//	s = append(s, src...)
//	return s
//}
//
//func calCapacity(c, l int) (int, bool) {
//	// 容量 <=64 缩不缩都无所谓，因为浪费内存也浪费不了多少
//	// 你可以考虑调大这个阈值，或者调小这个阈值
//	if c <= 64 {
//		return c, false
//	}
//	// 如果容量大于 2048，但是元素不足一半，
//	// 降低为 0.625，也就是 5/8
//	// 也就是比一半多一点，和正向扩容的 1.25 倍相呼应
//	if c > 2048 && (c/l >= 2) {
//		factor := 0.625
//		return int(float32(c) * float32(factor)), true
//	}
//	// 如果在 2048 以内，并且元素不足 1/4，那么直接缩减为一半
//	if c <= 2048 && (c/l >= 4) {
//		return c / 2, true
//	}
//	// 整个实现的核心是希望在后续少触发扩容的前提下，一次性释放尽可能多的内存
//	return c, false
//}
