package main

import "os"

func main() {
	puzzle := "maximum-width-of-binary-tree"
	question, err := fetch(puzzle)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("question.md", []byte(template(question)), 0644)
}
