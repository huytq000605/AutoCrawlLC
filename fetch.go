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

type question struct {
	Id         string `json:"questionFrontendId"`
	Content    string `json:"content"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Topics     []struct {
		Name string `json:"name"`
	} `json:"topicTags"`
	Hints []string `json:"hints"`
}

type leetcodeResponse struct {
	Data struct {
		Question question `json:"question"`
	} `json:"data"`
}

func fetchQuestion(puzzle, cookie string) (*question, error) {
	leetcode := "https://leetcode.com/graphql"
	query := []byte(fmt.Sprintf(`
	{"operationName":"questionData","variables":{"titleSlug":"%s"},"query":"query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n  title\n  content\n  difficulty\n questionFrontendId\n   topicTags { name\n }\n   hints\n }\n}\n"}
	`, puzzle))
	req, err := http.NewRequest("POST", leetcode, bytes.NewBuffer(query))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
  req.Header.Add("cookie", cookie)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responseParsed := leetcodeResponse{}
	err = json.Unmarshal(body, &responseParsed)

	if err != nil {
		return nil, err
	}
	return &responseParsed.Data.Question, nil
}

func fetchContest(contest, cookie string) ([]*question, error) {
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
	questionChan := make(chan question)
	var wg sync.WaitGroup
	questions := make([]*question, 0)

	for _, puzzle := range puzzles {
		wg.Add(1)
		go func(puzzle string) {
			defer wg.Done()
			question, err := fetchQuestion(puzzle, cookie)
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
