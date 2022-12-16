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
		puzzle := parts[3]

		question, err := fetchQuestion(puzzle)
		if err != nil {
			panic(err)
		}

		err = handleQuestion(question)
		if err != nil {
			panic(err)
		}

	case "contest":
		contest := parts[3]

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
