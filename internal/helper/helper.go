package helper

import "fmt"

func PrintError(err error) {
	PrintErrorString(err.Error())
}

func PrintErrorString(err string) {
	fmt.Println("Error occured: " + err)
}
