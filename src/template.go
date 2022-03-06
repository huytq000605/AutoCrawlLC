package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func template(question *questionType) string {
	return fmt.Sprintf(`
# %s. %s<br> %s

%s

<details>

<summary> Related Topics </summary>

%s

</details>

%s`,
		question.id, question.title, question.difficulty, getContent(question.title, question.content), getTopics(question.topics), getHints(question.hints))
}

func getContent(title, content string) string {
	r, err := regexp.Compile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	if err != nil {
		panic(err)
	}
	urls := r.FindAllString(content, -1)
	files, err := downloadAllFiles(title, urls)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(urls); i++ {
		content = strings.Replace(content, urls[i], fmt.Sprintf("./assets/%s", files[i]), 1)
	}

	return content
}

func getTopics(topics []string) string {
	if len(topics) == 0 {
		return ""
	}
	result := ""
	for _, topic := range topics {
		result += "\n"
		result += fmt.Sprintf("-\t`%s`", topic)
	}
	return result[1:]
}

func getHints(hints []string) string {
	if len(hints) == 0 {
		return ""
	}
	result := ""
	for i, hint := range hints {
		result += "\n\n"
		result += fmt.Sprintf("<details>\n<summary> Hint %d </summary>\n%s\n</details>", i+1, hint)
	}
	return result[1:]
}

func downloadAllFiles(title string, urls []string) ([]string, error) {
	if len(urls) == 0 {
		return []string{}, nil
	}
	result := make([]string, len(urls))
	assetsDir := filepath.Join(title, "assets")
	os.Mkdir(assetsDir, 07777)

	for idx, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		fileName := fmt.Sprintf("image%d%s", idx+1, filepath.Ext(url))
		result[idx] = fileName

		var file *os.File
		if _, err := os.Stat(filepath.Join(assetsDir, fileName)); err == nil {
			return nil, errors.New("File is already exists")
		} else {
			file, err = os.Create(filepath.Join(assetsDir, fileName))
			if err != nil {
				return nil, err
			}
		}

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
