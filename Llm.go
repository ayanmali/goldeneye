/*
Go implementation of Claude logic from Anthropic API
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LLM struct {
	Model  string
	ApiKey string
}

type LLMRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Tools       []Tool    `json:"tools"`
	Messages    []Message `json:"messages"`
	System      string    `json:"system"`
	Temperature float32   `json:"temperature"`
}

// Defines the structure of a single message (i.e. either system or user prompt)
type Message struct {
	Role    string `json:"role"`    // 'user' or 'assistant'
	Content string `json:"content"` // prompt
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"` // LLM response itself
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

/*
Gets the LLM's response from the user's request
*/
func (response LLMResponse) getOutput() string {
	var sb strings.Builder
	for _, content := range response.Content {
		if content.Type == "text" {
			sb.WriteString(content.Text)
			sb.WriteString("\n\n")
		}
	}
	return sb.String()
}

/*
Create a request for the LLM (custom messages parameter)
*/
func (llm LLM) NewLLMRequest(tools []Tool, maxTokens int, system string, temperature float32, messages []Message) *LLMRequest {
	return &LLMRequest{
		Model:       llm.Model,
		MaxTokens:   maxTokens,
		Tools:       tools,
		Messages:    messages,
		System:      system,
		Temperature: temperature,
	}
}

/*
Create a request that takes in a user prompt
*/
func (llm LLM) NewPromptRequest(tools []Tool, model string, maxTokens int, system string, temperature float32, prompt string) *LLMRequest {
	return &LLMRequest{
		Model:     llm.Model,
		MaxTokens: maxTokens,
		Tools:     tools,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		System:      system,
		Temperature: temperature,
	}
}

/*
Create a request that involves a user prompt as well as a partially filled response
*/
func (llm LLM) NewPromptRequestWithResponseStarter(tools []Tool, model string, maxTokens int, system string, temperature float32, prompt string, responseStarter string) *LLMRequest {
	return &LLMRequest{
		Model:     llm.Model,
		MaxTokens: maxTokens,
		Messages: []Message{
			{Role: "user", Content: prompt},
			{Role: "assistant", Content: responseStarter},
		},
		System:      system,
		Temperature: temperature,
	}
}

/*
Function for calling an LLM via a prompt.
Returns the Response, the status code of the request, and error, if applicable
*/
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
	req.Header.Set("x-api-key:", llm.ApiKey)
	req.Header.Set("anthropic-version:", "2024-10-22")
	req.Header.Set("content-type:", "application/json")

	// Sending the response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request %v\n", err)
		return nil, res.StatusCode, err
	}

	// Closing the response body once the function is finished executing
	defer res.Body.Close()

	// Reading the response data
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response %v\n", err)
		return nil, res.StatusCode, err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code %d: %s\n", res.StatusCode, string(resBody))
	}

	// Storing the response into the Response struct
	var apiRes LLMResponse
	err = json.Unmarshal(resBody, &apiRes)
	if err != nil {
		fmt.Printf("Error parsing the response: %v\n", err)
		return nil, res.StatusCode, err
	}

	return &apiRes, res.StatusCode, nil

}

//
/* TOOL USE LOGIC */
//

type Property struct {
	Type        string `json:"type"`        // the data type of this property (parameter)
	Description string `json:"description"` // a brief description of what this property (parameter) represents
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"` // properties (parameters) that this tool takes as input
	Required   []string            `json:"required"`   // names of properties that are mandatory input parameters
}

type Tool struct {
	Name        string      `json:"name"`         // i.e. name of the function/API endpoint to call
	Description string      `json:"description"`  // a brief description of what the tool does
	InputSchema InputSchema `json:"input_schema"` // defining what the tool uses as input
}

type SearchTool struct {
	apiKey      string // Brave API Key
	name        string
	description string
}
