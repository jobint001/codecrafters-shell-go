package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		var input string
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		input = line[:len(line)-1] // Remove the newline character
		if input == "exit" {
			return
		}
		command := strings.Split(input, " ")[0]
		switch command {
		case "echo":
			fmt.Println(input[5:]) // Print everything after "echo "
		default:	
		}
		fmt.Printf("%v: command not found\n", input)

	}

}
