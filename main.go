package main

import (
	"fmt"
)

func main() {

	llm := LLM{
		ApiKey: ANTHROPIC_API_KEY,
		Model:  "claude-3-haiku-20240307",
	}

	request := LLMRequest{
		Model:     llm.Model,
		MaxTokens: 1024,
		Messages: []Message{
			{Role: "user", Content: "How long is the Golden Gate Bridge?"},
		},
	}

	response, statusCode, err := llm.call(request)

	if statusCode < 200 || statusCode > 299 {
		fmt.Println("Error with HTTP request to LLM")
		return
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println(response.getOutput())

}
