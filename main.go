package main

import "fmt"

func main() {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7}
	//a1, err := Delete[int](arr, 3)
	//if err != nil {
	//	return
	//}
	//fmt.Println(a1)
	//fmt.Println(a1)
	a2, err := Delete[int](arr, 5)
	if err != nil {
		return
	}
	fmt.Println(a2)
}
