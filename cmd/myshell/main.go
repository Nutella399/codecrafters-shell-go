package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		repl(reader)
	}
}

func isBuiltin(command string) bool {
	value, exists := os.LookupEnv("PATH")
	if !exists {
		return false
	}

	switch command {
	case
		"echo",
		"exit",
		"type":
		return true
	default:
		pathArr := strings.Split(value, ":")
		fmt.Println(pathArr)
	}
	return false
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
		if isBuiltin(args[0]) {
			fmt.Println(args[0] + " is a shell builtin")
		} else {
			fmt.Println(args[0] + ": not found")
		}
	default:
		fmt.Println(text + ": command not found")
	}
}
