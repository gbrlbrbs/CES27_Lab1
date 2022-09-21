package utils

import (
	"fmt"
)

func PrintHelloWorld() {
	fmt.Println("Hello World!")
}

func CheckError(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
		panic(e)
	}
}