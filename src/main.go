package main

import (
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("Please input link")
	}
	parts := strings.Split(args[0], "/")
	removeBlankString(&parts)

	if parts[1] != "leetcode.com" {
		panic("Please input url from leetcode.com")
	}

	switch parts[2] {
	case "problems":
		var puzzle string
		for puzzle == "" || puzzle == "submissions" {
			puzzle = parts[len(parts)-1]
			parts = parts[:len(parts)-1]
		}

		question, err := fetchQuestion(puzzle)
		if err != nil {
			panic(err)
		}

		err = handleQuestion(question)
		if err != nil {
			panic(err)
		}

	case "contest":
		contest := parts[len(parts)-1]

		questions, err := fetchContest(contest)
		if err != nil {
			panic(err)
		}

		err = handleQuestions(questions)
		if err != nil {
			panic(err)
		}
	}

}

func removeBlankString(slice *[]string) {
	result := make([]string, 0, len(*slice))
	for _, str := range *slice {
		str = strings.Trim(str, " ")
		if str == "" {
			continue
		}
		result = append(result, str)
	}
	*slice = result
}
