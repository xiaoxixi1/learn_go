package main

import "fmt"

/*
*

				数组：语法 [cap]type
				   1 初始化要指定长度
			       2 可以直接初始化
		           3 len和cap用于获取数组的长度
	               4 使用for range循环
*/
func Array() {
	a1 := []int{9, 8, 7}
	fmt.Printf("a1: %v,len: %d,cap: %d", a1, len(a1), cap(a1))
}

// 没有初始化的元素默认填充0
func Array1() {
	a1 := [3]int{9, 8}
	fmt.Printf("a1: %v,len: %d,cap: %d", a1, len(a1), cap(a1))
	for _, val := range a1 {
		fmt.Print(val)
	}
}
