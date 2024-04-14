package main

import "fmt"

// go没有构造函数
func main() {
	// 结构体初始化
	u1 := &User{}
	println(u1) //是一个指针地址
	fmt.Printf("%+v \n", u1)

	u2 := User{}
	fmt.Printf("%+v \n", u2)

	//new可以理解为会为你的变量分配内存，并且把内存都置为0
	u3 := new(User)
	fmt.Printf("%+v \n", u3)

	var u4 *User
	fmt.Printf("%+v \n", u4)

	var u5 User
	fmt.Printf("%+v \n", u5)

	u6 := User{name: "Tom"}
	fmt.Printf("%+v \n", u6)

	u7 := User{"Jerry", 18}
	fmt.Printf("%+v \n", u7)
	println("==========================")
	ChangeUser()
	UseFish()
	Compose()
}
