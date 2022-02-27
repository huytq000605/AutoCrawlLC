package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func leetcode_question() url.Values {
	return url.Values{
		"operationName": {"questionData"},
	}
}

func main() {
	leetcodeUrl := "https://leetcode.com/graphql"
	js := []byte(`
	{"operationName":"questionData","variables":{"titleSlug":"maximum-width-of-binary-tree"},"query":"query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    questionFrontendId\n    boundTopicId\n    title\n    titleSlug\n    content\n    translatedTitle\n    translatedContent\n    isPaidOnly\n    difficulty\n    likes\n    dislikes\n    isLiked\n    similarQuestions\n    exampleTestcases\n    categoryTitle\n    contributors {\n      username\n      profileUrl\n      avatarUrl\n      __typename\n    }\n    topicTags {\n      name\n      slug\n      translatedName\n      __typename\n    }\n    companyTagStats\n    codeSnippets {\n      lang\n      langSlug\n      code\n      __typename\n    }\n    stats\n    hints\n    solution {\n      id\n      canSeeDetail\n      paidOnly\n      hasVideoSolution\n      paidOnlyVideo\n      __typename\n    }\n    status\n    sampleTestCase\n    metaData\n    judgerAvailable\n    judgeType\n    mysqlSchemas\n    enableRunCode\n    enableTestMode\n    enableDebugger\n    envInfo\n    libraryUrl\n    adminUrl\n    challengeQuestion {\n      id\n      date\n      incompleteChallengeCount\n      streakCount\n      type\n      __typename\n    }\n    __typename\n  }\n}\n"}
	`)
	req, err := http.NewRequest("POST", leetcodeUrl, bytes.NewBuffer(js))
	req.Header.Add("referer", "https://leetcode.com/problems/maximum-width-of-binary-tree/")
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
	fmt.Println(string(body))
}
