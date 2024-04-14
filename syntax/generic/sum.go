package main

import (
	"io"
)

/**
  泛型方法
*/

// 求和
func Sum[T Number](vals []T) T {
	var res T
	for _, val := range vals {
		res = res + val // 如果T直接约束any，会报错，因为不是所有的类型都可以用加法
	}
	return res
}

// 找到最大值
func Max[T Number](vals []T) T {
	max := vals[0]
	for _, val := range vals {
		if max < val {
			max = val
		}
	}
	return max
}

// 找到最小值
func Min[T Number](vals []T) T {
	min := vals[0]
	for _, val := range vals {
		if min > val {
			min = val
		}
	}
	return min
}

// 过滤与查找
func find[T Number](vals []T, filter func(t T) bool) T {
	for _, val := range vals {
		if filter(val) {
			return val
		}
	}
	var t T
	return t
}

// 在指定位置插入
func Insert[T Number](index int, t T, vals []T) []T {
	if index < 0 || index > len(vals) {
		return nil
	}
	vals = append(vals, t)
	for i := len(vals) - 1; i > index; i-- {
		vals[i] = vals[i-1]
	}
	vals[index] = t
	return vals
}

// 当使用Number约束，number中约束数值类型，比如int，uint等，则上面使用+不会报错
type Number interface {
	~int | uint | int32 // ~int表示int的衍生类型，原类型支持的话，衍生类型也可以支持
}

// 除了可以约束成自定义的接口，还可以约束成普通的接口
func Closable[T io.Closer]() {
	var t T
	/**	err := t.Close()
	if err != nil {
		return
	}*/
	t.Close() // t这里是有close方法，因为T必须实现close方法
}

type Integer int

func userSum() {
	list := []int{123, 234}
	res := Sum[int](list)
	println("sum:", res)

	list1 := []Integer{234, 234}
	res1 := Sum[Integer](list1)
	println("sum:", res1)
}
