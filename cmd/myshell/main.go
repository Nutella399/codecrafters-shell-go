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

func repl(reader *bufio.Reader) {
	fmt.Fprint(os.Stdout, "$ ")

	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	textArr := strings.Split(text, " ")

	switch textArr[0] {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(textArr[1:], " "))
	default:
		fmt.Println(text + ": command not found")
	}
}
