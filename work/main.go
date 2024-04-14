package main

import "fmt"

func main() {
	//场景1 : cap 小于等于256，直接删除不需要缩容
	slice := make([]int, 256, 256)
	slice[2] = 99
	slice = deleteByIndex[int](2, slice)
	fmt.Printf("cap: %d , len: %d , slice[2]: %d \n", cap(slice), len(slice), slice[2])
	// 场景2 : cap大于256，但是(cap - len)*2 < cap,直接删除，不需要缩容
	slice1 := make([]int, 259, 259)
	slice1[2] = 69
	slice1 = deleteByIndex[int](2, slice1)
	fmt.Printf("cap: %d , len: %d , slice[2]: %d \n", cap(slice1), len(slice1), slice1[2])
	// 场景3： cap大于256 ，并且cap - len)*2 > cap,删除，需要缩容
	slice3 := make([]int, 151, 300)
	slice3[2] = 59
	slice3 = deleteByIndex[int](2, slice3)
	fmt.Printf("cap: %d , len: %d , slice[2]: %d \n", cap(slice3), len(slice3), slice3[2])

}
