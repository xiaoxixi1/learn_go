package main

import (
	"fmt"
	"io"
)

/*
*

			组合：
			   接口组合接口
		       结构体组合结构体
		       结构体组合结构体指针
		       结构体组合接口
	        可以组合多个
*/
type Inner struct {
	//	name string
}

// 结构体组合结构体
type Outer struct {
	Inner
}

// 结构体组合结构体指针,区别在于，在使用inner的方法是需要对Inner进行初始化
type Outer1 struct {
	*Inner
}

// 组合接口
type Outer2 struct {
	io.Closer
}

func (i Inner) Name() string {
	return "Inner"
}

func (o Outer) Name() string {
	return "Outer"
}

func (i Inner) Hello() {
	fmt.Println("Hello,I am ", i.Name())
}

/*
*

			组合特性：
			  1 A组合B之后，可以在A上调用所有B的方法
		      2 B实现的所有接口都认为A已经实现了
	          3 组合不是继承，没有多态
*/
func Compose() {
	out := Outer{}
	out.Hello()
	// 需要对inner初始化
	out1 := Outer1{&Inner{}}
	out1.Hello()
	// 其实就是将inner当成out的成员变量
	out2 := Outer1{}
	out2.Inner = &Inner{}
	out2.Hello()

	out3 := Outer2{}
	fmt.Println(out3)

	out4 := Outer{}
	out4.Hello()
}
