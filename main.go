package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	_cookiePath = "./cookie_lc"
)

func parseFlag() (needLogin bool) {
	flag.BoolVar(&needLogin, "l", false, "Need login flag")
	flag.Visit(func(f *flag.Flag) {
		fmt.Println(f)
		if f.Name == "l" || f.Name == "login" {
			needLogin = true
		}
	})
	flag.Parse()
	return needLogin
}

func main() {
	needLogin := parseFlag()
	if needLogin {
		if err := ExtractCookies(); err != nil {
			panic(err)
		}
	}
	cookieJson, err := os.ReadFile(_cookiePath)
	if err != nil {
		fmt.Println("No cookie found")
	}

	var cookies []*Cookie
	if err := json.Unmarshal(cookieJson, &cookies); err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		panic("Please input link")
	}
	parts := strings.Split(args[len(args)-1], "/")
	if len(parts) < 4 {
		panic("The URL is not correct")
	}
	removeBlankString(&parts)

	if parts[1] != "leetcode.com" {
		panic("Please input url from leetcode.com")
	}

	switch parts[2] {
	case "problems":
		puzzle := parts[3]
		question, err := fetchQuestion(puzzle, cookies)
		if err != nil {
			panic(err)
		}

		err = handleQuestion(question)
		if err != nil {
			panic(err)
		}

	case "contest":
		contest := parts[3]

		questions, err := fetchContest(contest, cookies)
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
