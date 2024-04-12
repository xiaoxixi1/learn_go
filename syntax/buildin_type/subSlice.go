package main

import "fmt"

// 数组和切片都可以通过[start:end]的形式来获取子切片,左闭右开
// 子切片和原切片是共享底层数组的
func Subslice() {
	slice1 := []int{1, 2, 3, 4, 5}
	sub := slice1[0:2]
	fmt.Printf("a1: %v,len: %d,cap: %d \n", sub, len(sub), cap(sub))
	sub1 := slice1[:1]
	fmt.Printf("a1: %v,len: %d,cap: %d \n", sub1, len(sub1), cap(sub1))
	sub2 := slice1[1:]
	fmt.Printf("a1: %v,len: %d,cap: %d \n", sub2, len(sub2), cap(sub2))
}

// 子切片和切片去修改数组可能会互相影响，只有当切片的数据结构发生变化，才不会共享底层数组
// 当切片发生扩容时，会生成一个新的数组，然后再将原来数组里面的数据复制过去，此时，子切片和切片不再共享底层数组
func Shareslice() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	fmt.Printf("a1: %v,len: %d,cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("a1: %v,len: %d,cap: %d \n", s2, len(s2), cap(s2))
	s2[0] = 99
	fmt.Printf("a1: %v,len: %d,cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("a1: %v,len: %d,cap: %d \n", s2, len(s2), cap(s2))
	s2 = append(s2, 199)
	s2[0] = 1999
	fmt.Printf("a1: %v,len: %d,cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("a1: %v,len: %d,cap: %d \n", s2, len(s2), cap(s2))

}
