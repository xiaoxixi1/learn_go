package wire

import "fmt"

func UserRepository() {
	repo := InitUserRepository() // 这里使用的生成的wire_gen.go里面的方法
	fmt.Println(repo)
}
