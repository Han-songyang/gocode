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
	return arr
}
