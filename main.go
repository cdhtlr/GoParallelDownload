package main

import (
	"GoParallelDownload/function/download"
	"fmt"
)

func main() {
	completed, err := download.Run() //returns bool and error, can be used in conditional statements
	if completed {
		fmt.Println("Download Completed :D")
	} else {
		fmt.Println(err)
	}
}
