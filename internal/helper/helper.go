package helper

import "fmt"

// PrintError is a wrapper of PrintErrorString
func PrintError(err error) {
	PrintErrorString(err.Error())
}

// PrintErrorString prints error messages on stdout
func PrintErrorString(err string) {
	fmt.Println("Error occured: " + err)
}
