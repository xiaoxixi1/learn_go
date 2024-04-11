package main

import (
	"math"
	"project_go/syntax/variables/demo"
	"strconv"
	"strings"
	"unicode/utf8"
)

/*
*

	基础类型：数值类型：uint,int,float；字符串类型；byte类型，bool类型
	工具包：数值型math，字符串类型strings,byte类型bytes
	注意点：
*/
func main() {
	println("hello,Go!")
	Hello()
	var a int = 345
	var b int = 123
	a--
	a++
	// float默认是float64
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
	if math.Abs(c-d) < math.SmallestNonzeroFloat64 {
		println("c=d")
	} else {
		println("c != d")
	}

	String()
	Byte()
	Bool()
	// 当问全局变量
	println(demo.D)
}

func String() {
	println("hello" + "world")
	var a = 123
	println("hello", string(a))
	println("hello" + string(rune(1)))
	println(strings.Split("a+b", "+")[0])
	// 双引号可以使用反斜杠转义，但是反引号不行
	println("hello \"")
	// 反引号可以换行
	println(`hello 
换行了`)
	// 拼接字符串和别的类型
	println("hello" + strconv.Itoa(123))
	// len计算的是字节数，不是字符数
	println(len("hello你好"))
	println(utf8.RuneCountInString("hello你好"))
}
func Byte() {
	// byte和string可以互相转换，但是byte但是使用一般很少，一般使用的是byte的切片[]byte
	var str string = "abc"
	var bt []byte = []byte(str)
	var str1 string = string(bt)
	println(str1)
	// print打印的是byte字符的ascall码
	var a = 'a'
	var c = 12
	println(a)
	println(c)
}

func Bool() {
	var a bool = true
	var b bool = false
	println(a && b)
	println(a || b)
	println(!a)
	// !(a&&b) = !a || !b
	// !(a||b) = !a && !b

}
