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

//auxiliary function to determine the bigger between two numbers
func MaxNumber(x1, x2 int) int {
	if x1 > x2 {
		return x1
	}
	return x2
}

//auxiliary function to concatenate ID, clock and 
func Concatenate(str_id string, str_clock string, text string) (message string) {
	message =  str_id + "," + str_clock + "," + text

	return
}