package main

import (
	"fmt"
	"os"
	"sync"
)

func handleQuestion(question *questionType) error {
	err := os.Mkdir(question.title, 07777)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s/question.md", question.title), []byte(template(question)), 0644)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Crawled %s successfuly", question.title))
	return nil
}

func handleQuestions(questions []*questionType) error {
	// errChan := make(chan error)
	var wg sync.WaitGroup
	for _, question := range questions {
		wg.Add(1)
		go func(question *questionType) {
			defer wg.Done()
			err := handleQuestion(question)
			if err != nil {
				// errChan <- err
				fmt.Println(err)
				return
			}
		}(question)
	}
	wg.Wait()
	return nil
}
