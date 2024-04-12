package main

// go里面if里面可以定义局部变量
func ifVariable() {
	end := 10
	start := 5
	if distance := end - start; distance > 4 {
		println(distance)
	}
	// 在if语句里面定义的局部变量，在if语句外不能引用
	//println(distance)
}
