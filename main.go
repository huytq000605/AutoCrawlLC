package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type questionType struct {
	content    string
	title      string
	difficulty string
	topics     []string
	hints      []string
}

func main() {
	leetcode := "https://leetcode.com/graphql"
	puzzle := "sliding-window-maximum"
	query := []byte(fmt.Sprintf(`
	{"operationName":"questionData","variables":{"titleSlug":"%s"},"query":"query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n  title\n  content\n  difficulty\n questionFrontendId\n   topicTags { name\n }\n   hints\n }\n}\n"}
	`, puzzle))
	req, err := http.NewRequest("POST", leetcode, bytes.NewBuffer(query))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var responseMap map[string]interface{}
	json.Unmarshal(body, &responseMap)

	questionInResponse := responseMap["data"].(map[string]interface{})["question"].(map[string]interface{})

	question := questionType{}

	question.content = questionInResponse["content"].(string)
	question.difficulty = questionInResponse["difficulty"].(string)
	question.title = questionInResponse["questionFrontendId"].(string) + ". " + questionInResponse["title"].(string)

	topicTags := questionInResponse["topicTags"].([]interface{})
	question.topics = make([]string, len(topicTags))
	for idx, topicTag := range topicTags {
		question.topics[idx] = topicTag.(map[string]interface{})["name"].(string)
	}

	hints := questionInResponse["hints"].([]interface{})
	question.hints = make([]string, len(hints))
	for idx, hint := range hints {
		question.hints[idx] = hint.(string)
	}

	err = os.WriteFile("test.md", []byte(template(question)), 0644)
	if err != nil {
		panic(err)
	}

	// fmt.Println(template(question))

}
