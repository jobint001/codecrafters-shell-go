package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
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
			continue
		case "type":
			handleTypeCommand(input[5:])
			continue
		default:
		}
		fmt.Printf("%v: command not found\n", input)

	}

}

func handleTypeCommand(input string) {
	
	switch input {
	case "echo", "exit", "type":
		fmt.Printf("%v is a shell builtin\n", input)

	default:
		path := os.Getenv("PATH")
		paths := strings.Split(path, ":")
		fmt.Println(paths)
		for _,p := range paths {
			path, err := exec.LookPath(fmt.Sprintf("%s/%s", p, input))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v is %s",input,path)
			return
			//fmt.Printf(path)
		}

		fmt.Printf("%v: not found\n", input)

	}

}
