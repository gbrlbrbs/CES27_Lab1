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
func maxNumber(x1, x2 int) int {
	if x1 > x2 {
		return x1
	}
	return x2
}

//auxiliary function to concatenate ID, clock and 
func concatenate(str_id string, str_clock string, text string) string {
	message :=  str_id + "," + str_clock + "," + text
	buf := []byte(message)

	return message
}

// auxiliary non-blocking async routine to listen for terminal input
func readInput(ch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch <- string(text)
	}
}