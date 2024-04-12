package main

/*
*
for语句表现形式：

		1 for 三段式
		2 类似while
		3 死循环
	    4 遍历数组和切片，map
*/
func for1() {
	for i := 0; i < 2; i++ {
		println(i)
	}
}
func for2() {
	i := 0
	for i < 2 {
		println(i)
	}
}

func for3() {
	i := 0
	for {
		if i > 10 {
			break
		}
		println("")
		i++
	}
}

func forArray() {
	arr := [3]string{"A", "B", "C"}
	for index, value := range arr {
		print(index)
		println(value)
	}
}
func forSlice() {
	// 切片和数组的区别就是[]中没有长度
	arr := []string{"A", "B", "C"}
	for index, value := range arr {
		print(index)
		println(value)
	}
	// 忽略值
	for inx := range arr {
		println(inx, arr[inx])
	}
	// 忽略下标
	for _, val := range arr {
		println(val)
	}
}

func forMap() {
	m := map[string]string{
		"1": "A",
		"2": "B",
	}
	for key, value := range m {
		println(key)
		println(value)
	}
	for key := range m {
		println(m[key])
	}
	for _, val := range m {
		println(val)
	}

}

type user struct {
	name string
}

func LoopBug() {
	users := []*user{
		{
			name: "Tom",
		},
		{
			name: "Jerry",
		},
	}
	m := make(map[string]*user)
	for _, value := range users {
		m[value.name] = value
	}
	for key, value := range m {
		println(key, value.name)
	}
}
