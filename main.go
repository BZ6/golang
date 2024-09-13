package main

import "fmt"

func createHello() (result *string) {
	str := "Hello GO"
	result = &str
	return
}

func main() {
	hello := createHello()
	fmt.Println(*hello)
}
