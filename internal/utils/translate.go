package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"sync"

	"miniflux.app/v2/internal/model"
)

// chatgpt response struct
type Message struct {
	Role    string
	Content string
}
type Choice struct {
	Finish_reason string
	Index         int
	Message       Message
}

type Usage struct {
	Adjust_total      int32
	Completion_tokens int32
	Final_total       int
	Pre_token_count   int32
	Pre_total         int
	Prompt_tokens     int
	Total_tokens      int
}

type ChatGPTResponse struct {
	Id                 string
	Model              string
	Object             string
	System_fingerprint string
	Choices            []Choice
	Created            string
	Usage              Usage
}

func TranslateByGPT(entry *model.Entry, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	msg := entry.Title
	url := "https://oa.api2d.net/v1/chat/completions"
	key := ""
	prompt := "你是一个翻译助手，需要将我的话准确地翻译成中文"
	model := "gpt-3.5-turbo"

	if key == "" || prompt == "" {
		return
	}
	method := "POST"
	data := fmt.Sprintf(`{
		"model": "%s",
		"messages": [
			{
				"role": "system",
				"content": "%s"
			},
			{
				"role": "user",
				"content": "%s"
			}
		],
		"safe_mode": false
	}`, model, prompt, msg)
	payload := strings.NewReader(data)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(body))
	var response ChatGPTResponse
	jsonErr := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", jsonErr)
		return
	} else {
		choices := response.Choices
		if len(choices) != 0 {
			content := choices[0].Message.Content
			(*entry).Title += "｜" + content
			// fmt.Printf("chat gpt response:%v\n", response)
		}
		return
	}
}

func PostProcessEntriesTitle(feed *model.Feed, entries *model.Entries) {
	if !feed.Translatable {
		return
	}
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}
	var wg sync.WaitGroup
	for _, entry := range *entries {
		wg.Add(1)
		go TranslateByGPT(entry, client, &wg)
	}
	wg.Wait()
}
