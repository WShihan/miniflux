package translate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"log/slog"
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
	System_fingerprint string
	Choices            []Choice
	Usage              Usage
}

type ChatGPT struct {
	Key   string
	To    string
	Model string
	URL   string
}

func (chatgpt *ChatGPT) GetKey() string {
	return chatgpt.Key
}

func (chatgpt *ChatGPT) Execute(sem chan struct{}, entry *model.Entry, client *http.Client, wg *sync.WaitGroup, ak string) {
	defer wg.Done()

	prompt := fmt.Sprintf("You are an assistant, please translate my word into %s", chatgpt.To)

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
	}`, chatgpt.Model, prompt, entry.Title)
	payload := strings.NewReader(data)

	req, err := http.NewRequest(method, chatgpt.URL, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ak))
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
	var response ChatGPTResponse
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		fmt.Println("Error parsing JSON:", jsonErr)
		slog.Error(fmt.Sprintf("Translate error: %s", entry.Title))
		return
	} else {
		choices := response.Choices
		if len(choices) != 0 {
			content := choices[0].Message.Content
			slog.Info(fmt.Sprintf("Translate title:%s, result:%s", entry.Title, content))
			entry.Title += "｜" + content
		}
		return
	}
}
