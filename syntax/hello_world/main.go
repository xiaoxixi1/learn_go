package main

import "math"

func main() {
	println("hello,Go!")
	Hello()
	var a int = 345
	var b int = 123
	a--
	a++
	c := 234.1
	d := 678.0
	c++
	println(c)
	println(a + b)
	println(c / d)
	//只有同类型的数值才能进行加减乘除，go中没有自动类型转换
	println(float64(a) * c)
	// 计算使用math包
	println(math.Abs(-123.4))
}
