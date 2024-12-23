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

func isInPath(command string) (string, bool) {
	value, exists := os.LookupEnv("PATH")
	if !exists {
		return "", false
	}

	pathArr := strings.Split(value, ":")
	for _, path := range pathArr {
		filename := path + "/" + command
		_, err := os.Stat(filename)
		if !errors.Is(err, os.ErrNotExist) {
			return filename, true
		}
	}

	return "", false
}

func parseCommand(command string) []string {
	var result []string
	var temp strings.Builder
	singleQuoted := false
	doubleQuoted := false
	backSlashed := false

	for i := 0; i < len(command); i++ {
		char := command[i]

		switch char {
		case ' ':
			if singleQuoted || doubleQuoted || backSlashed {
				temp.WriteByte(char)
				if backSlashed {
					backSlashed = false
				}
			} else if temp.Len() > 0 {
				result = append(result, temp.String())
				temp.Reset()
			}
		case '\'':
			if doubleQuoted {
				temp.WriteByte(char)
			} else if singleQuoted {
				singleQuoted = false
			} else {
				singleQuoted = true
			}
		case '"':
			if doubleQuoted {
				doubleQuoted = false
			} else if singleQuoted {
				temp.WriteByte(char)
			} else {
				doubleQuoted = true
			}
		case '\\':
			if !singleQuoted && (i < len(command)-1 && command[i+1] == '\\' || command[i+1] == '$' || command[i+1] == '"') {
				nextChar := command[i+1]
				temp.WriteByte(nextChar)
				i++
			} else if singleQuoted || doubleQuoted {
				temp.WriteByte(char)
			} else {
				backSlashed = true
			}
		default:
			temp.WriteByte(char)
		}
	}

	if temp.Len() > 0 {
		result = append(result, temp.String())
	}
	return result
}

func repl(reader *bufio.Reader) {
	fmt.Fprint(os.Stdout, "$ ")

	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	textArr := parseCommand(text)

	command := textArr[0]
	args := textArr[1:]

	switch command {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "type":
		filename, existsInPath := isInPath(args[0])
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
		path := args[0]
		if args[0] == "~" {
			path = os.Getenv("HOME")
		}
		err := os.Chdir(path)
		if err != nil {
			fmt.Println("cd: " + args[0] + ": No such file or directory")
		}
	default:
		_, existsInPath := isInPath(command)
		if existsInPath {
			cmd := exec.Command(command, args...)
			stdout, err := cmd.Output()
			if err != nil {
				return
			}
			fmt.Print(string(stdout))
		} else {
			fmt.Println(text + ": command not found")
		}
	}
}
