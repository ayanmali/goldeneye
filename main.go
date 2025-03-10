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
func getWeather(location string) float32 {
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

func getWeatherVariadic(args ...any) any {
	if len(args) == 0 {
		return "Error: No location provided"
	}

	location, ok := args[0].(string)
	if !ok {
		return "Error: Invalid location type"
	}

	return getWeather(location)
}

func main() {
	hasSystemPrompt := false
	hasMemory := true

	// Replace w actual values
	systemPrompt := ""
	memoryBlock := MemoryBlock{}

	llm := LLM{
		ApiKey: ANTHROPIC_API_KEY,
		Model:  "claude-3-haiku-20240307",
		// System: []Content{
		// 	{Type: "text", Text: "You are a helpful AI assistant designed to answer user questions."},
		// },
		Messages: []Message{
			{Role: "user",
				Content: []Content{
					{Type: "text", Text: "What is the weather in New York City?"},
				}},
		},
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
				Function: getWeatherVariadic,
			},
		},
	}

	if hasSystemPrompt {
		llm.System = append(llm.System, Content{Type: "text", Text: systemPrompt})
	}
	if hasMemory {
		llm.MemoryBlock = memoryBlock
		llm.System = append(llm.System, Content{Type: "text", Text: ""})
		llm.Tools = append(llm.Tools,
			// Core memory save Tool
			Tool{
				Name:        "coreMemorySave",
				Description: "Save important information about you (the agent) or the human you are chatting with.",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]Property{
						"agentMem": {
							Type:        "boolean",
							Description: "Must be either `true` to save information about yourself (the agent), or `false` to save information about the user.",
						},
						"memory": {
							Type:        "string",
							Description: "Memory to save in the section",
						},
					},
					Required: []string{"agentMem", "memory"},
				},
				// Wrapper function
				Function: memoryBlock.coreMemorySaveVariadic,
				// Function: func(agentMema any, memorya any) any {
				// 	//args := params.(map[string]any)
				// 	agentMem := params["agentMem"].(bool)
				// 	memory := params["memory"].(string)
				// 	return memoryBlock.coreMemorySave(agentMem, memory)
				// },
			},
		)
		llm.System = append(llm.System, Content{Type: "text", Text: fmt.Sprintf("MEMORY:\n%+v", memoryBlock)}) // provide string representation of MemoryBlock struct
	}

	request := LLMRequest{
		Model:  llm.Model,
		System: llm.System,
		// System: []Content{
		// 	{Type: "text", Text: "You are a helpful AI assistant designed to answer user questions."},
		// },

		Messages: llm.Messages,

		MaxTokens: 1024,
		Tools:     llm.Tools,
	}

	response, statusCode, err := llm.call(request)

	if statusCode < 200 || statusCode >= 300 {
		fmt.Println("Error with HTTP request to LLM")
		return
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// LLM's initial output
	fmt.Println(response.getOutput())
	fmt.Println("---------")

	// Adding LLM response to chat history
	llm.addResponseToChatHistory(*response)

	// Checking each Content block returned in the response
	for _, content := range response.Content {
		// only looking for tool_use Content blocks
		if content.Type != "tool_use" {
			continue
		}
		// Adding tool call result to chat history
		llm.addToolResultToChatHistory(content)
	}

	fmt.Println(llm.Messages)

	response1, statusCode1, err1 := llm.call(request)

	if statusCode1 < 200 || statusCode >= 300 {
		fmt.Println("Error with HTTP request to LLM")
		return
	}

	if err1 != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println(response1.getOutput())

}
