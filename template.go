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

func template(question *question) string {
	return fmt.Sprintf(`
# %s. %s<br> %s

%s

<details>

<summary> Related Topics </summary>

%s

</details>

%s`,
		question.Id, question.Title, question.Difficulty, getContent(question.Title, question.Content), getTopics(question), getHints(question.Hints))
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

	for i := 0; i < len(files); i++ {
		if files[i] == urls[i] {
			continue
		}
		content = strings.Replace(content, urls[i], fmt.Sprintf("./assets/%s", files[i]), 1)
	}

	return content
}

func getTopics(q *question) string {
	if len(q.Topics) == 0 {
		return ""
	}
	result := ""
	for _, topic := range q.Topics {
		result += "\n"
		result += fmt.Sprintf("-\t`%s`", topic.Name)
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
	fileNames := make([]string, len(urls))
	assetsDir := filepath.Join(title, "assets")
	os.Mkdir(assetsDir, 07777)

	ignore := 0

	for idx, url := range urls {
		acceptedExtensions := []string{".apng", ".avif", ".gif", ".jpg", ".jpeg", ".jfif", ".pjpeg", ".pjp", ".png"}
		extension := filepath.Ext(url)

		if include(acceptedExtensions, extension) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			fileName := fmt.Sprintf("image%d%s", idx+1, extension)
			fileNames[idx] = fileName

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
		} else {
			fileNames[idx] = urls[idx]
			ignore += 1
		}
	}
	if ignore == len(urls) {
		os.Remove(assetsDir)
	}
	return fileNames, nil
}

func include(slice []string, s string) bool {
	for _, str := range slice {
		if s == str {
			return true
		}
	}
	return false
}
