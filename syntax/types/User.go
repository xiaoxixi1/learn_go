package main

import "fmt"

type User struct {
	name string
	age  int
	//user *User //结构体自引用，只能使用指针
	// 准确的来说，如果在整个引用链上构成了循环，就只能使用指针
}

/*
*

	方法接收器
	一个方法既可以定义在结构体上，也可以定义在这个结构体的指针上
	 前者：结构体接收器  后者：指针接收器
*/

/**
  结构体接收器修改字段，内部不会生效，指针接收器修改才会生效
  因为方法调用本身是值传递。值传递的意思就是发生复制
  如果是基本类型或者结构体，相当于复制了一份
  如果是指针，那么就复制一份指针，但是指向的结构体还是同一个
  内置类型复制了一份，但是内部数据结构还是同一个
*/

func (u User) ChangeName(name string) {
	u.name = name
	fmt.Printf("%+p \n", &u)
}

/*
*
相当于
*/
func ChangeName(u User, name string) {
	u.name = name
	fmt.Printf("%+p \n", &u)
}

func (u *User) ChangeAge(age int) {
	u.age = age
	fmt.Printf("%+p \n", u)
}

/*
*
相当于
*/
func ChangeAge(u *User, age int) {
	u.age = age
	fmt.Printf("%+p \n", u)
}

func ChangeUser() {
	u1 := User{name: "Tome", age: 18}
	fmt.Printf("%+v \n", u1)
	// ChangeName也会打印u1的地址
	u1.ChangeName("Jerry")
	// ChangeName也会打印u1的地址
	u1.ChangeAge(16)
	// 打印u1的地址，会发现调用ChangeName时，实际时复制的地址，修改的复制的u1的name,所以打印u1会发现u1的name没有变化
	fmt.Printf("%+p \n", &u1)
	// 但是age会产生变化，因为age传入的时指针
	fmt.Printf("%+v \n", u1)

	println("=============u2==============")
	// 就是u2定义是一个指针地址，结果也是一样的
	u2 := &User{name: "Tom", age: 18}
	fmt.Printf("%+v \n", u2)
	u2.ChangeName("Jerry") // 这里之所以不变，是因为这里传进去的还是结构。复制的也是结构体
	//ChangeName(u2,"Jerry")// 会报错，因为ChangeName这个函数的参数是结构体，而不是指针
	u2.ChangeAge(16)
	fmt.Printf("%+p \n", u2)
	// 但是age会产生变化，因为age传入的时指针
	fmt.Printf("%+v \n", u2)

}
