package main

import (
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("Please input link to question")
	}
	parts := strings.Split(args[0], "/")
	var puzzle string
	for puzzle == "" {
		puzzle = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}

	question, err := fetch(puzzle)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(question.title, 07777)
	if err != nil {
		panic(err)
	}

	err = os.Chdir(question.title)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("question.md", []byte(template(question)), 0644)
}
