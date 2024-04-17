package main

import (
	"golang.org/x/crypto/bcrypt"
)

/**
  使用bcrypt加密算法
  bcrypt号称目前最安全的加密算法，不需要传入盐值
  可以通过cost来控制加密性能
  同样的文本。两次加密的结果不同
  bcrypt加密密码长度要小于72
*/

func main() {
	password := []byte("123456#hello")
	// 加密
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	println(err != nil)
	// 没有解密，只有判断密码是否一致
	println(bcrypt.CompareHashAndPassword(encrypted, password) == nil)
}
