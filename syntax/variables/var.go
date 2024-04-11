package main

import "project_go/syntax/variables/demo"

/**
  变量类型：局部变量，块变量，全局变量
*/

// 块声明，其中的变量为包变量，作用范围在包内，也就是当前文件夹内
var (
	b = 234
	c = "hello"
)

// 也是属于包变量，包变量只声明没有使用不会编译报错
var d = "world"

func main() {
	// 局部变量
	var a = 123
	println(a)
	// 全局变量，包外可访问
	println(demo.D)
	// 但是局部变量不能只声明不使用
	//var f =123
	// :=变量声明的方式只能用于局部变量
	g := "sdf"
	println(g)
	// 不同作用域，变量可以申明多次
	var c = 123
	println(c)
}
