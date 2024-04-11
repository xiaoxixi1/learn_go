package demo

// 首字母大写，表示包外也可以使用的全局变量
var D = "Hi"

// iota 在 const关键字出现时将被重置为 0(const 内部的第一行之前)，const 中每新增一行常量声明将使 iota 计数一次(iota 可理解为 const 语句块中的行索引)。
const (
	a = iota
	b
	c
)

// iota常用来初始化常量
const (
	mystatus = iota<<1 + 1
	mystatus1
	mystatus2
)
