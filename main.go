/*
In Context
Core CoreMemory
-> to determine when to search out of context, part of core memory is a Statistics section, which contains an overview of the kinds of data, number of entries, etc. present in both kinds of out of context memory
-> Agent can examine the Statistics section to determine whether to search OOC memory

Out of Context
Recall - previous chat history
Archival - general information (Ex. code, PDFs, etc)
- requires semantic search in a vector DB to retrieve relevant data

To clear space in the context window and send data from core memory to recall/archival:
Flush
Recursive Summary
*/
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
	memory := CoreMemory{}

	llm := Agent{
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
		llm.CoreMemory = memory
		llm.System = append(llm.System, Content{Type: "text", Text: "You are a helpful AI agent with access to a core memory. You can add new information to your core memory by calling the `coreMemorySave` function."})
		llm.Tools = append(llm.Tools,
			// Core memory save Tool
			Tool{
				Name:        "coreMemoryAppend",
				Description: "Save important information about you (the agent) or the human you are chatting with.",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]Property{
						"agentMem": {
							Type:        "boolean",
							Description: "Must be either `true` to save information about yourself (the agent), or `false` to save information about the user.",
						},
						"newMemory": {
							Type:        "string",
							Description: "CoreMemory to save in the section",
						},
					},
					Required: []string{"agentMem", "newMemory"},
				},
				// Wrapper function
				Function: llm.coreMemoryAppendVariadic,
				// Function: func(agentMema any, memorya any) any {
				// 	//args := params.(map[string]any)
				// 	agentMem := params["agentMem"].(bool)
				// 	memory := params["memory"].(string)
				// 	return memoryBlock.coreMemorySave(agentMem, memory)
				// },
			},
		)
		llm.System = append(llm.System, Content{Type: "text", Text: fmt.Sprintf("MEMORY:\n%+v", memory)}) // provide string representation of MemoryBlock struct
	}

	request := AgentRequest{
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
		fmt.Println("Error with HTTP request to Agent")
		return
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Agent's initial output
	fmt.Println(response.getOutput())
	fmt.Println("---------")

	// Adding Agent response to chat history
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
		fmt.Println("Error with HTTP request to Agent")
		return
	}

	if err1 != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println(response1.getOutput())

}
