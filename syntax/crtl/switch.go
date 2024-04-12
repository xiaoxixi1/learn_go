package main

import "fmt"

func Switch(status string) {
	switch status {
	case "1":
		fmt.Printf("", status)
	case "2":
		fmt.Printf("", status)
	default:
		fmt.Printf("", "hello")
	}
}

func Switch2(status int) {
	switch {
	case status >= 18:
		fmt.Printf("", "成年了")
	case status < 6:
		fmt.Printf("", "婴幼儿")
	default:
		println("default")
	}
}
