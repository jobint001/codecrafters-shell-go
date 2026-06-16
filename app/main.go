package main

import (
	"fmt"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	 fmt.Print("$ ")
	 var i string
	 _, err := fmt.Scanln(&i)
	 if err != nil {
		 fmt.Println("Error reading input:", err)
		 return
	 }
	fmt.Printf("%v: command not found",i)
}
