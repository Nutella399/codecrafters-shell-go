package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		repl(reader)
	}
}

func isBuiltin(command string) bool {
	switch command {
	case
		"echo",
		"exit",
		"type",
		"pwd",
		"cd":
		return true
	default:
		return false
	}
}

func isInPath(command string) (bool, string) {
	value, exists := os.LookupEnv("PATH")
	if !exists {
		return false, ""
	}

	pathArr := strings.Split(value, ":")
	for _, path := range pathArr {
		filename := path + "/" + command
		_, err := os.Stat(filename)
		if !errors.Is(err, os.ErrNotExist) {
			return true, filename
		}
	}

	return false, ""
}

func repl(reader *bufio.Reader) {
	fmt.Fprint(os.Stdout, "$ ")

	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	textArr := strings.Split(text, " ")

	command := textArr[0]
	args := textArr[1:]

	switch command {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "type":
		existsInPath, filename := isInPath(args[0])
		if isBuiltin(args[0]) {
			fmt.Println(args[0] + " is a shell builtin")
		} else if existsInPath {
			fmt.Println(filename)
		} else {
			fmt.Println(args[0] + ": not found")
		}
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error printing directory: ", err)
		} else {
			fmt.Println(pwd)
		}
	case "cd":
		err := os.Chdir(args[0])
		if err != nil {
			fmt.Println("cd: " + args[0] + ": No such file or directory")
		}
	default:
		existsInPath, _ := isInPath(command)
		if existsInPath {
			cmd := exec.Command(command, args[0])
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Print(string(stdout))
		} else {
			fmt.Println(text + ": command not found")
		}
	}
}
