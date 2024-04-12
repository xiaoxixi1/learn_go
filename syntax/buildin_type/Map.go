package main

// 初始化
func Map() {
	// 直接初始化
	m1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	println(m1)
	// make初始化
	m2 := make(map[string]string, 4)
	m2["key2"] = "value2"
	println(m2)

	/**
	  map读取，有两个返回值
	  返回值1：value
	  返回值2：这个元素是否存在
	  如果只返回一个元素，则默认是value
	*/
	val, ok := m1["key1"]
	println(val)
	println(ok)
	println(m2["key2"])
	println(len(m1))
	// map每次遍历的顺序都是随机的
	for k, v := range m1 {
		println(k, v)
	}
	for k := range m1 {
		println(k, m1[k])
	}
	for _, v := range m1 {
		println(v)
	}

	//删除元素delete
	delete(m1, "key1")

}

/**
  comparable（可比较的）概念：在编译和运行的时候能够判断元素是否是相等的
  在switch里面和在map里面的key，都必须是可比较的
     基本类型和string都是可比较的
     如果元素是可比较的，那对应的数组也是可比较的
  另外go是强类型语言，无法进行类型转换
*/
