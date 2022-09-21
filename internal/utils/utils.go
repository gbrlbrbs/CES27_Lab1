package utils

import (
	"fmt"
)

func PrintError(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
	}
}

func CheckError(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
		panic(e)
	}
}