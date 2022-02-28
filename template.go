package main

import "fmt"

func template(question questionType) string {
	return fmt.Sprintf(`
# %s<br> %s

%s

<details>

<summary> Related Topics </summary>

%s

</details>

%s
	`, question.title, question.difficulty, question.content, getTopics(question.topics), getHints(question.hints))
}
func getTopics(topics []string) string {
	result := ""
	for _, topic := range topics {
		result += "\n"
		result += fmt.Sprintf("-\t`%s`", topic)
	}
	return result[1:]
}

func getHints(hints []string) string {
	result := ""
	for i, hint := range hints {
		result += "\n\n"
		result += fmt.Sprintf("<details>\n<summary> Hint %d </summary>\n%s\n</details>", i+1, hint)
	}
	return result[1:]
}
