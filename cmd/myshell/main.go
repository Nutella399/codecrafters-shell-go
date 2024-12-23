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

func writeToFile(file string, output string) error {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		writeErr := os.WriteFile(file, []byte(output), 0o777)
		return writeErr
	} else {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			return readErr
		}
		newContent := append(content, []byte(output)...)
		writeErr := os.WriteFile(file, newContent, 0o777)
		return writeErr
	}
}

func repl(reader *bufio.Reader) {
	fmt.Fprint(os.Stdout, "$ ")

	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\r\n")
	textArr := parseCommand(text)

	stdOutput := ""
	stdError := ""
	command := textArr[0]
	args := textArr[1:]
	var secondaryCommand []string

	for index, arg := range args {
		if arg == "1>" || arg == "1>>" || arg == ">" || arg == ">>" || arg == "2>" || arg == "2>>" {
			args = textArr[1 : index+1]
			secondaryCommand = textArr[index+1:]
		}
	}

	switch command {
	case "exit":
		os.Exit(0)
	case "echo":
		stdOutput = strings.Join(args, " ") + "\n"
	case "type":
		filename, existsInPath := isInPath(args[0])
		if isBuiltin(args[0]) {
			stdOutput = args[0] + " is a shell builtin\n"
		} else if existsInPath {
			stdOutput = filename + "\n"
		} else {
			stdOutput = args[0] + ": not found\n"
		}
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			stdError = "Error printing directory: " + err.Error() + "\n"
		} else {
			stdOutput = pwd + "\n"
		}
	case "cd":
		path := args[0]
		if args[0] == "~" {
			path = os.Getenv("HOME")
		}
		err := os.Chdir(path)
		if err != nil {
			stdOutput = "cd: " + args[0] + ": No such file or directory\n"
		}
	default:
		_, existsInPath := isInPath(command)
		if existsInPath {
			cmd := exec.Command(command, args...)
			stdout, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					stdError = string(exitErr.Stderr)
				} else {
					stdError = "error:" + err.Error()
				}
			}
			stdOutput = string(stdout)
		} else {
			stdOutput = text + ": command not found\n"
		}
	}

	if len(secondaryCommand) > 0 {
		secondCommand := secondaryCommand[0]
		secondArgs := secondaryCommand[1:]

		switch secondCommand {
		case ">", ">>", "1>", "1>>":
			err := writeToFile(secondArgs[0], stdOutput)
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Fprint(os.Stderr, stdError)
		case "2>", "2>>":
			err := writeToFile(secondArgs[0], stdError)
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Fprint(os.Stdout, stdOutput)
		}
	} else {
		if len(stdOutput) > 0 {
			fmt.Fprint(os.Stdout, stdOutput)
		}
		if len(stdError) > 0 {
			fmt.Fprint(os.Stderr, stdError)
		}
	}
}
