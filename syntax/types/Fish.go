package main

import (
	"fmt"
	"go/types"
)

type Fish struct {
}

func (f1 Fish) Swim() {

}

/*
*

				衍生类型：
				   语法：type TypeA TypeB
				   衍生类型，是一种新的类型
			       使用场景：想扩展第三方库结构体的方法
		           衍生类型和原类型可以互相转换， a :=TypeA(typeB类型的变量) b:=TypeB(TypeA类型的变量)
	               TypeB实现了某个接口不等于TypeA也实现了某个接口
*/
type FishFake Fish

func (f FishFake) FakeSwim() {

}

type StringFake types.Map

/*
*

	类型别名：
	  语法：type TypeA = TypeB
*/
type Yu = Fish

func UseFish() {
	f1 := Fish{}
	f1.Swim()
	f2 := FishFake{}
	// f2.Swim() 会报错
	f2.FakeSwim()
	f3 := Fish(f2)
	f4 := FishFake(f1)
	fmt.Println(f3, f4)

	s1 := StringFake{}
	fmt.Println(s1)
	f5 := Yu{}
	fmt.Println(f5)
}
