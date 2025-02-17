package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
curl https://api.anthropic.com/v1/messages \
     --header "x-api-key: $ANTHROPIC_API_KEY" \
     --header "anthropic-version: 2023-06-01" \
     --header "content-type: application/json" \
     --data \
'{
    "model": "claude-3-5-sonnet-20241022",
    "max_tokens": 1024,
    "messages": [
        {"role": "user", "content": "Hello, world"}
    ]
}'
*/

type LLM struct {
	model  string
	apiKey string
}

type LLMRequest struct {
	model     string    `json:"model"`
	maxTokens int       `json:"max_tokens"`
	messages  []Message `json:"messages"`
}

// Defines the structure of a single message (i.e. either system or user prompt)
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResponseMessage struct {
	Content []Content `json:"content"`
	Role    string    `json:"role"`
}

type LLMResponse struct {
	Content      []Content       `json:"content"`
	ID           string          `json:"id"`
	Model        string          `json:"model"`
	Role         string          `json:"role"`
	Type         string          `json:"type"`
	StopReason   string          `json:"stop_reason"`
	StopSequence any             `json:"stop_sequence"`
	Usage        map[string]int  `json:"usage"`
	Message      ResponseMessage `json:"message"`
}

func (llm LLM) call(reqData LLMRequest) (*LLMResponse, int, error) {

	// Extracting the request data
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		fmt.Printf("Error marshaling request data: %v\n", err)
		return nil, -1, err
	}

	// Creating the request
	url := "https://api.anthropic.com/v1/messages"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Printf("Error creating request:%v\n", err)
		return nil, -1, err
	}

	// Setting request headers
	req.Header.Set("x-api-key:", llm.apiKey)
	req.Header.Set("anthropic-version:", "2024-10-22")
	req.Header.Set("content-type:", "application/json")

	// Sending the response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request %v\n", err)
		return nil, -1, err
	}

	// Closing the response body once the function is finished executing
	defer res.Body.Close()

	// Reading the response data
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response %v\n", err)
	}

	// Storing the response into the Response struct
	var apiRes LLMResponse
	err = json.Unmarshal(resBody, &apiRes)
	if err != nil {
		fmt.Printf("Error parsing the response: %v\n", err)
		return nil, -1, err
	}

	return &apiRes, res.StatusCode, nil

}
