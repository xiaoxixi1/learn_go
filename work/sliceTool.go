package main

/**
实现删除切片特定下标元素的方法。
要求一：能够实现删除操作就可以。
要求二：考虑使用比较高性能的实现。
要求三：改造为泛型方法
要求四：支持缩容，并旦设计缩容机制。
扩容机制：在容量为256之前，每次扩容2倍，在容量为256之后，每次扩容1.25倍
*/

/*
*

			删除指定下标元素：
			  index: 待删除元素下标
			  vals: []泛型切片
			 return: 删除指定元素后的切片
		    缩容规则：
	            提供两个参数；缩容阈值，缩容因子。这里默认256和1.25
	               当分片的长度小于等于这个阈值，则随意删除，不进行缩容；
	               当分片的长度大于这个阈值时，如果cap-len=cap/2 ,则根据缩容因子缩容
	           一般建议根据业务，传入合适的缩容阈值和缩容因子
*/
const sliceShinkThreshold = 256
const shrinkFactor float32 = 1.25

func deleteByIndex[T any](index int, vals []T) []T {
	if index < 0 || index > len(vals) {
		return vals
	}
	cap := cap(vals)
	len := len(vals)
	// case 1: 切片容量小于256，直接删除元素，不缩容
	if cap <= sliceShinkThreshold {
		// 直接删除
		return deleteByIndexDerict(index, len, vals)
	}
	// 判断删除后是否需要缩容
	if (cap-len+1)*2 < cap {
		return deleteByIndexDerict(index, len, vals)
	}
	shinkedLen := int(float32(cap) / shrinkFactor)
	shinkedSlice := make([]T, len-1, shinkedLen)
	j := 0
	for i := 0; i < len; i++ {
		if i == index {
			continue
		}
		shinkedSlice[j] = vals[i]
		j++
	}
	return shinkedSlice
}

func deleteByIndexDerict[T any](index int, len int, vals []T) []T {
	for i := index; i < len-1; i++ {
		vals[i] = vals[i+1]
	}
	return vals[:len-1]
}
