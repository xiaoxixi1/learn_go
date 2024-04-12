package main

import "fmt"

/*
*

			切片：语法[]type
			初始化方式：直接初始化，make初始化
			访问：与数组一致
		    切片可以类比Arraylist去理解，但是对比ArrayList的增删改查，只有append
	        切片支持子切片
*/
func Slice() {
	slice1 := []int{4, 2, 3}
	fmt.Printf("a1: %v,len: %d,cap: %d \n", slice1, len(slice1), cap(slice1))

	//len是当前切片元素的个数，cap是当前切片的容量
	slice2 := make([]int, 3, 4)
	fmt.Printf("a1: %v,len: %d,cap: %d \n", slice2, len(slice2), cap(slice2))

	slice2 = append(slice2, 5)
	fmt.Printf("a1: %v,len: %d,cap: %d \n", slice2, len(slice2), cap(slice2))

	//自动扩容（耗性能）
	slice2 = append(slice2, 6)
	fmt.Printf("a1: %v,len: %d,cap: %d \n", slice2, len(slice2), cap(slice2))

}
