package main

/*
*
泛型结构体
*/
type List[T any] interface {
	ADD(index int, T any)
	Append(T any)
	Delete(val int)
}

type LinkList[T any] struct {
	head *NodeV1[T]
}

type NodeV1[T any] struct {
	data T
}

func (l LinkList[T]) Add(index int, val T) {

}

func userList() {
	list := LinkList[int]{}
	list.Add(1, 123)
	// list.Add(2,"123") 会报错，存在类型约束

}
