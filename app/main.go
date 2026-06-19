package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/google/shlex"
)

var builtInCommands = []string{
	"echo",
	"type",
	"exit",
	"pwd",
	"cd",
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		fields, _ := shlex.Split(command)

		if len(fields[0]) == 0 {
			os.Exit(0)
		}
		args := fields[1:]
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			os.Exit(1)
		}

		if slices.Contains(builtInCommands, fields[0]) {
			switch fields[0] {
			case "exit":
				os.Exit(0)
			case "echo":
				handleEcho(args)
			case "type":
				handleType(fields[1])
			case "pwd":
				handlePwd()
			case "cd":
				handleCd(fields[1])
			}
		} else if _, err := exec.LookPath(fields[0]); err == nil {
			cmd := exec.Command(fields[0], fields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			fmt.Printf("%s: command not found\n", fields[0])
		}

	}
}

func handleType(command string) {
	if slices.Contains(builtInCommands, command) {
		fmt.Printf("%s is a shell builtin\n", command)
		return
	} else if path, err := exec.LookPath(command); err == nil {
		fmt.Printf("%s is %s\n", command, path)
		return
	}

	fmt.Printf("%s: not found\n", command)
}

func handlePwd() {
	path, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println(path)

}
func handleCd(input string) {

	if input == "~" {
		homePath := os.Getenv("HOME")
		os.Chdir(homePath)
		return
	}

	_, err := os.Stat(input)

	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("cd: %s: No such file or directory\n", input)
		return
	}
	os.Chdir(input)

}

func handleEcho(args []string) {
	var output strings.Builder
	var file string
	if slices.Contains(args, ">") {
		for i, value := range args {
			if args[i] == ">" {
				file = args[i+1]
				break
			} else {
				output.WriteString(value)
			}

		}
		
		f, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = f.WriteString(output.String())
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}

	} else {
		fmt.Println(strings.Join(args, " "))
	}

}
