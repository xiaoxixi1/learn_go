package main

import "fmt"

// 方法签名：方法名+参数+返回值

func main() {
	Invoke()
	Func5("黎晨曦")
	Func5("黎晨曦", "晨曦", "xi")
	Func6()
	Defer()
	DeferClosure()
	DeferClosure1()
	DeferClosure2()
	println(DeferReturn())
	println(DeferReturn2())
	println(DeferReturn3().name)
	DeferClosureLoop1()
	DeferClosureLoop2()
	DeferClosureLoop3()
}

func Invoke() {
	Func1()
	Func2("我调用Func2了")
	println(Func3("我调用Func3了"))
	str, err := Func3("")
	println(str, err)
	_, err = Func4("")
	println(err)
	MyFunc4()
	Myfunc5()
	fn := Myfunc6()
	println(fn("world"))
	Myfunc7()

}

// 不带参数，不带返回值的函数
func Func1() {
	println("func1")
}

// 带参数不带返回的函数
func Func2(str string) {
	println("func2:" + str)
}

// 带参数带返回的函数
func Func3(str string) (string, error) {
	println("func3:", str)
	return "hello,world", nil
}

// 带参数带名字的返回的函数
func Func4(str string) (str1 string, err error) {
	return
}

/*
*

		不定参数：alise可以传入int最大值个参数，也可以不传
		不定参数一定是放再最后一个参数
	    不定参数在方法内部可被当成切片使用
*/
func Func5(name string, alise ...string) {
	if len(alise) != 0 {
		println(alise[0])
	}
}

func Func6(alise ...int) {
	fmt.Println("不定参数，我不传任何值")
}
