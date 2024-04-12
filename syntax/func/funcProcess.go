package main

/*
*

	函数式编成入门：函数作为变量进行赋值
*/
func Myfunc3() (string, error) {
	println("hello,myFunc3")
	return "", nil
}

func MyFunc4() {
	myFunc := Myfunc3
	_, _ = myFunc()
}

/*
*
函数式编程：局部方法，在方法内部声明一个局部方法，作用域就在当前方法内
*/
func Myfunc5() {
	fn := func(name string) string {
		return "hello" + name
	}
	str := fn("Bob")
	println(str)
}

/*
*
函数式编程：方法本身作为返回值
*/
func Myfunc6() func(name string) string {
	return func(name string) string {
		return "hello" + name
	}
}

/*
*
函数式编程：匿名方法发起调用
声明一个方法后，发起调用
*/
func Myfunc7() {
	func(name string) {
		println("hello," + name)
	}("myfunc7")
	fn := func(name string) string {
		return "hello," + name
	}("myfunc77")
	println(fn)
	c := Closure("大明")
	println(c())
}

/**
函数式编程：闭包
闭包：方法+它绑定运行时上下文
一个对象如果被闭包引用，不会进行垃圾回收，所以使用不当可能会导致内存泄露
*/
