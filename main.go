package main

import (
	"fmt"
)

// func getWeather(location string) float32 {
//     if location == "San Francisco, CA" {
//         return 24.0
//     } else if location == "Boston, MA" {
//         return 10.0
//     } else if location == "New York City, NY" {
//         return 8.0
// 		}
//	  return 5.0
// }

/*
TODO:
Add logic for redefining a function in the variadic parameter format
*/
func getWeather(args ...any) any {
	if len(args) == 0 {
		return "Error: No location provided"
	}

	location, ok := args[0].(string)
	if !ok {
		return "Error: Invalid location type"
	}

	switch location {
	case "San Francisco, CA":
		return float32(24.0)
	case "Boston, MA":
		return float32(10.0)
	case "New York City, NY":
		return float32(8.0)
	default:
		return float32(5.0)
	}
}

func main() {

	llm := LLM{
		ApiKey: ANTHROPIC_API_KEY,
		Model:  "claude-3-haiku-20240307",
	}

	request := LLMRequest{
		Model: llm.Model,
		// System: []Content{
		// 	{Type: "text", Text: "You are a helpful AI assistant designed to answer user questions."},
		// },

		Messages: []Message{
			{Role: "user",
				Content: []Content{
					{Type: "text", Text: "What is the weather in New York City?"},
				}},
		},

		MaxTokens: 1024,
		Tools: []Tool{
			{Name: "getWeather",
				Description: "A function that returns the weather (in degrees Celsius) for a given location.",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]Property{
						"location": {Type: "string", Description: `The geographic location to check the weather for. 
																	The format should be the name of the city followed by the abbreviated name of the state that the city belongs to.
																	In other words, the format should be "<CITY_NAME>, <STATE>".
																	For example, to check the weather in Boston, Massachussetts, this parameter would be "Boston, MA".`,
						},
					},
					Required: []string{"location"},
				},
				Function: getWeather,
			},
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

	fmt.Println("---------")

	// Adding LLM response to chat history
	request.addResponseToChatHistory(*response)

	// Checking each Content block returned in the response
	for _, content := range response.Content {
		// only looking for tool_use Content blocks
		if content.Type != "tool_use" {
			continue
		}
		// Adding tool call result to chat history
		request.addToolResultToChatHistory(content)
	}

	fmt.Println(request.Messages)

	response1, statusCode1, err1 := llm.call(request)

	if statusCode1 < 200 || statusCode > 299 {
		fmt.Println("Error with HTTP request to LLM")
		return
	}

	if err1 != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println(response1.getOutput())

}
