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

		if err != nil {
			// EOF (e.g. Ctrl-D) terminates the shell.
			os.Exit(0)
		}
		if len(fields) == 0 {
			// Empty line: just show a fresh prompt.
			continue
		}

		fields, redirectFile, redirectErr, append, rerr := parseRedirect(fields)
		if rerr != nil {
			fmt.Fprintln(os.Stderr, rerr)
			continue
		}

		// Resolve where stdout should go.
		stdout := os.Stdout
		var f *os.File
		var ferr error
		if redirectFile != "" {
			if append {
				f, ferr = os.Open(redirectFile)
				if ferr != nil {
					fmt.Fprintln(os.Stderr,ferr)
					goto middle
				}
				defer f.Close()
				goto end
			}
			middle:
			f, ferr = os.Create(redirectFile)
			if ferr != nil {
				fmt.Fprintln(os.Stderr, ferr)
				continue
			}
			end:
			stdout = f
		}

		stderr := os.Stderr
		if redirectErr != "" {
			f, ferr := os.Create(redirectErr)
			if ferr != nil {
				fmt.Fprintln(os.Stderr, ferr)
				continue
			}
			stderr = f
		}
		args := fields[1:]

		if slices.Contains(builtInCommands, fields[0]) {
			switch fields[0] {
			case "exit":
				os.Exit(0)
			case "echo":
				fmt.Fprintln(stdout, strings.Join(args, " "))
			case "type":
				handleType(fields[1])
			case "pwd":
				handlePwd()
			case "cd":
				handleCd(fields[1])
			}
		} else if _, err := exec.LookPath(fields[0]); err == nil {
			cmd := exec.Command(fields[0], fields[1:]...)
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			cmd.Run()
		} else {
			fmt.Printf("%s: command not found\n", fields[0])
		}

		if stdout != os.Stdout {
			stdout.Close()
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

// parseRedirect scans fields for a stdout redirection operator (">" or "1>")
// and returns the command fields with the operator + target removed, plus the
// target filename ("" if there is no redirection).
func parseRedirect(fields []string) ([]string, string, string, bool, error) {
	for i, f := range fields {
		if f == ">" || f == "1>" {
			if i+1 >= len(fields) {
				return nil, "", "", false, fmt.Errorf("syntax error: expected filename after %q", f)
			}
			return fields[:i], fields[i+1], "", false, nil
		} else if f == "2>" {
			return fields[:i], "", fields[i+1], false, nil

		} else if f == ">>" || f=="1>>" {
			return fields[:i], "", fields[i+1], true, nil
		}
	}
	return fields, "", "", false, nil
}
