package main

// 接口只能定义函数,是一组行为的抽象
type List interface {
	ADD(index int, val string)
	Append(val string)
	Delete(val int)
}
type LinkList struct {
	/**
	  实现接口的快捷键：
	    在此处右键，选择generate,然后再输入要实现的接口，回车
	*/
	head string
}

func (l *LinkList) ADD(index int, val string) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkList) Append(val string) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkList) Delete(val int) {
	//TODO implement me
	panic("implement me")
}

// 当一个结构体具体接口所有的方法的时候，它就实现了这个接口
func (l *LinkList) Add(index int, val string) {
	// 实现ADD方法
}
