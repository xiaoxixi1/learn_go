package main

var i = "world"

func Closure(name string) func() string {
	// 闭包
	// 上下文：变量i和name
	// 方法：匿名函数func
	return func() string {
		return "hello," + name + i
	}
}
