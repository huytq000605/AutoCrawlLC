package main

import (
	"fmt"
	"os"
	"sync"
)

func handleQuestion(question *question) error {
	err := os.Mkdir(question.Title, 07777)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s/question.md", question.Title), []byte(template(question)), 0644)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Crawled %s successfuly", question.Title))
	return nil
}

func handleQuestions(questions []*question) error {
	// errChan := make(chan error)
	var wg sync.WaitGroup
	for _, q := range questions {
		wg.Add(1)
		go func(q *question) {
			defer wg.Done()
			err := handleQuestion(q)
			if err != nil {
				// errChan <- err
				fmt.Println(err)
				return
			}
		}(q)
	}
	wg.Wait()
	return nil
}
