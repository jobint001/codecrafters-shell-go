package main

import (
	"bufio"
	"fmt"
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
			if len(input) > len(command) {
				runExternalPgm(command, input[len(command)+1:])

			}

		}

	}

}

func handleTypeCommand(input string) {

	switch input {
	case "echo", "exit", "type":
		fmt.Printf("%v is a shell builtin\n", input)

	default:
		path := os.Getenv("PATH")
		paths := strings.Split(path, string(os.PathListSeparator))
		//fmt.Println(paths)
		for _, p := range paths {
			err := os.Chdir(p)
			path, err := exec.LookPath(fmt.Sprintf("%s/%s", p, input))
			if err != nil {
				continue
			}
			fmt.Printf("%v is %s\n", input, path)
			return
			//fmt.Printf(path)
		}

		fmt.Printf("%v: not found\n", input)

	}

}

func runExternalPgm(command, input string) {

	args := strings.Split(input, " ")
	path := os.Getenv("PATH")
	paths := strings.SplitSeq(path, string(os.PathListSeparator))
	//fmt.Println(paths)
	for p := range paths {
		err := os.Chdir(p)
		path, err = exec.LookPath(fmt.Sprintf("%s", command))
		if err != nil {
			continue
		}
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		return
		//fmt.Printf(path)
	}
	fmt.Printf("%v: command not found\n", input)

}
