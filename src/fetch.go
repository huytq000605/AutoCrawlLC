package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

type questionType struct {
	id         string
	content    string
	title      string
	difficulty string
	topics     []string
	hints      []string
}

func fetchQuestion(puzzle string) (*questionType, error) {
	leetcode := "https://leetcode.com/graphql"
	query := []byte(fmt.Sprintf(`
	{"operationName":"questionData","variables":{"titleSlug":"%s"},"query":"query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n  title\n  content\n  difficulty\n questionFrontendId\n   topicTags { name\n }\n   hints\n }\n}\n"}
	`, puzzle))
	req, err := http.NewRequest("POST", leetcode, bytes.NewBuffer(query))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseMap map[string]interface{}
	json.Unmarshal(body, &responseMap)

	questionInResponse := responseMap["data"].(map[string]interface{})["question"].(map[string]interface{})

	question := questionType{}

	question.content = questionInResponse["content"].(string)
	question.difficulty = questionInResponse["difficulty"].(string)
	question.title = questionInResponse["title"].(string)
	question.id = questionInResponse["questionFrontendId"].(string)

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

	if err != nil {
		return nil, err
	}
	return &question, nil
}

func fetchContest(contest string) ([]*questionType, error) {
	leetcode := fmt.Sprintf("https://leetcode.com/contest/api/info/%s/", contest)
	resp, err := http.Get(leetcode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseMap map[string]interface{}
	json.Unmarshal(body, &responseMap)

	questionMap := responseMap["questions"].([]interface{})
	puzzles := make([]string, 0)

	for _, question := range questionMap {
		puzzle := question.(map[string]interface{})["title_slug"].(string)
		puzzles = append(puzzles, puzzle)
	}

	doneChan := make(chan struct{})
	errChan := make(chan error)
	questionChan := make(chan questionType)
	var wg sync.WaitGroup
	questions := make([]*questionType, 0)

	for _, puzzle := range puzzles {
		wg.Add(1)
		go func(puzzle string) {
			defer wg.Done()
			question, err := fetchQuestion(puzzle)
			if err != nil {
				errChan <- err
				return
			}
			questionChan <- *question
		}(puzzle)
	}

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	for {
		select {
		case question := <-questionChan:
			questions = append(questions, &question)
		case <-doneChan:
			return questions, nil
		case err := <-errChan:
			return nil, err
		}
	}
}
