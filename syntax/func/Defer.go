package main

/**
go里面有一段进制：允许在方法返回前，执行一段逻辑 ==》defer
defer 也叫延迟调用
*/

// defer类似于栈，先进后出，后进先出
func Defer() {
	defer func() {
		println("调用第一个函数")
	}()

	defer func() {
		println("调用第二个函数")
	}()
}

/*
*
defer与闭包
defer值确认的原则：

	作为参数传入时：在定义defer时已经确认了
	作为闭包引入时：执行defer定义的方法时才确认
*/
func DeferClosure() {
	var j = 1
	defer func() {
		println(j)
	}()
}
func DeferClosure1() {
	var j = 1
	defer func() {
		println(j)
	}()
	j = 2
	j = 3
}

func DeferClosure2() {
	var j = 1
	defer func(i int) {
		println(i)
	}(j)
	j = 2

}

/*
*
defer修改返回值
如果是带名字的返回值，则可以进行修改，否则修改失效
*/
func DeferReturn() int {
	a := 1
	defer func() {
		a = 2
	}()
	return a // 1
}

func DeferReturn2() (a int) {
	a = 1
	defer func() {
		a = 2
	}()
	return a //2
}

func DeferReturn3() *Mystruct {
	res := &Mystruct{
		name: "Tom",
	}
	defer func() {
		res.name = "jerry"
	}()
	return res // name jerry 因为没有改返回值，而是改的res这个返回值指针指向的值
}

type Mystruct struct {
	name string
}

// 练习
func DeferClosureLoop1() {
	println("DeferClosureLoop1:")
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
	println("=============================")
}

func DeferClosureLoop2() {
	println("DeferClosureLoop2:")
	for i := 0; i < 10; i++ {
		defer func(i int) {
			println(i)
		}(i)
	}
	println("=============================")
}

// 这里的结果是9876543210
func DeferClosureLoop3() {
	println("DeferClosureLoop3:")
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
	println("=============================")
}
