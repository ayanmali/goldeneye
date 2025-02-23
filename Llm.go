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
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens,omitempty"`
	Tools     []Tool `json:"tools,omitempty"`
	/*
		Specify how the LLM should use the given tools. 'any' forces a tool to be used, 'tool' forces a specific tool to be used, and 'auto' does not enforce any tool usage.
		Follows schema: {"type" : "auto/any/tool"} --- If "type" is set to "tool", then you must include another key "name", which is the name of the tool to force the LLM to use.
		Parallel tool use can be disabled by adding `"disable_parallel_tool_use" = true`.
	*/
	ToolChoice map[string]string `json:"tool_choice,omitempty"`
	/*
		Defines the conversation history with the LLM
	*/
	Messages    []Message `json:"messages"`
	System      string    `json:"system,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

// Defines the structure of a single message (i.e. either system or user prompt)
type Message struct {
	Role string `json:"role"` // 'user' or 'assistant'
	/*
			Can either be a string or []Content.
			Example:
			messages=[
		        {
		            "role": "user",
		            "content": "What's the weather like in San Francisco?"
		        },
		        {
		            "role": "assistant",
		            "content": [
		                {
		                    "type": "text",
		                    "text": "<thinking>I need to use get_weather, and the user wants SF, which is likely San Francisco, CA.</thinking>"
		                },
		                {
		                    "type": "tool_use",
		                    "id": "toolu_01A09q90qw90lq917835lq9",
		                    "name": "get_weather",
		                    "input": {"location": "San Francisco, CA", "unit": "celsius"}
		                }
		            ]
		        }
			]
	*/
	Content any `json:"content"` // prompt
}

// type Source struct {
// 	// E.g. "base64"
// 	Type string `json:"type"`
// 	// E.g. "image/jpeg"
// 	MediaType string `json:"media_type"`
// 	// E.g. "/9j/4AAQSkZJRg..."
// 	Data string `json:"data"`
// }

type Content struct {
	/*
		The type of content being provided back to the user.
		The type can be 'text', 'tool_use' (used when providing a Tool to the LLM for it to use), or 'tool_result' (used when passing the output of a Tool call back to the LLM)
	*/
	Type  string            `json:"type"`
	Text  string            `json:"text,omitempty"` // LLM response itself
	ID    string            `json:"id,omitempty"`
	Name  string            `json:"name,omitempty"`
	Input map[string]string `json:"input,omitempty"`
	/*
		The ID from the initial tool use call to the LLM.
		This ID is generated in the response after making an initial LLM call with tool use enabled.
		Use this key with 'role' = 'user' when passing the result from a tool call back to the LLM. Be sure to also include the 'content' key
	*/
	ToolUseID string `json:"tool_use_id,omitempty"`
	/*
		The output from the tool call.
		Use this key when passing the result from the tool call back to the LLM
	*/
	Content string `json:"content,omitempty"`
}

// type ResponseMessage struct {
// 	Content []Content `json:"content"`
// 	Role    string    `json:"role"`
// }

type LLMResponse struct {
	/*
		Contains multiple Content structs which correspond to different parts of the LLM's output.
		E.g. type "text" is the LLM's raw response, type "tool_use" is the model's tool call.
	*/
	Content []Content `json:"content"`
	ID      string    `json:"id"`
	Model   string    `json:"model,omitempty"`
	Role    string    `json:"role,omitempty"`
	//Type         string          `json:"type"`
	StopReason   string         `json:"stop_reason,omitempty"`
	StopSequence any            `json:"stop_sequence,omitempty"`
	Usage        map[string]int `json:"usage,omitempty"`
	//Message      ResponseMessage `json:"message"`
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
func (llm LLM) NewPromptRequest(maxTokens int, system string, temperature float32, prompt string) *LLMRequest {
	return &LLMRequest{
		Model:     llm.Model,
		MaxTokens: maxTokens,
		//Tools:     tools,
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
func (llm LLM) NewPromptRequestWithResponseStarter(tools []Tool, maxTokens int, system string, temperature float32, prompt string, responseStarter string) *LLMRequest {
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
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Printf("Error creating request:%v\n", err)
		return nil, -1, err
	}

	// Setting request headers
	fmt.Println("Setting request headers")
	req.Header.Set("x-api-key", llm.ApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	fmt.Println("Making request")

	// Sending the response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request %v\n", err)
		return nil, res.StatusCode, err
	}

	// Closing the response body once the function is finished executing
	defer res.Body.Close()

	fmt.Println("Reading response data")
	// Reading the response data
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response %v\n", err)
		return nil, res.StatusCode, err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code %d: %s\n", res.StatusCode, string(resBody))
	}

	fmt.Println("Storing response")
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
	Type        string   `json:"type"`        // the data type of this property (parameter)
	Description string   `json:"description"` // a brief description of what this property (parameter) represents
	Enum        []string `json:"enum"`        // slice of possible values for Enum type input parameters
}

type InputSchema struct {
	Type string `json:"type"`
	/*
		Properties (parameters) that this tool takes as input.
		Keys are the names of the tools.
		Values are Property structs, which denote the type and description for the tool, as well as enumerated values, if applicable.
		Ensure that the Properties are listed in the order in which the function calls them.
		For example, if a function `f` takes in `a` and `b` as parameters and should be called in the form `f(a,b)`,
		then `a` should be listed before `b` in this field.
	*/
	Properties map[string]Property `json:"properties"` // properties (parameters) that this tool takes as input
	/*
		Names of properties that are mandatory input parameters.
		This slice should contain keys from Properties that are mandatory for being called
	*/
	Required []string `json:"required"`
}

/*
Provide extremely detailed descriptions.
This is by far the most important factor in tool performance.
Your descriptions should explain every detail about the tool, including:
- What the tool does
- When it should be used (and when it shouldn’t)
- What each parameter means and how it affects the tool’s behavior
Any important caveats or limitations, such as what information the tool does not return if the tool name is unclear.
The more context you can give Claude about your tools, the better it will be at deciding when and how to use them.
Aim for at least 3-4 sentences per tool description, more if the tool is complex.

While you can include examples of how to use a tool in its description or in the accompanying prompt,
this is less important than having a clear and comprehensive explanation of the tool’s purpose and parameters.
Only add examples after you’ve fully fleshed out the description.

Example of a good Tool description:

	{
	  "name": "get_stock_price",
	  "description": "Retrieves the current stock price for a given ticker symbol. The ticker symbol must be a valid symbol for a publicly traded company on a major US stock exchange like NYSE or NASDAQ. The tool will return the latest trade price in USD. It should be used when the user asks about the current or most recent price of a specific stock. It will not provide any other information about the stock or company.",
	  "input_schema": {
	    "type": "object",
	    "properties": {
	      "ticker": {
	        "type": "string",
	        "description": "The stock ticker symbol, e.g. AAPL for Apple Inc."
	      }
	    },
	    "required": ["ticker"]
	  }
	}

Example of a bad Tool description:

	{
	  "name": "get_stock_price",
	  "description": "Gets the stock price for a ticker.",
	  "input_schema": {
	    "type": "object",
	    "properties": {
	      "ticker": {
	        "type": "string"
	      }
	    },
	    "required": ["ticker"]
	  }
	}
*/
type Tool struct {
	Name        string      `json:"name"`         // i.e. name of the function/API endpoint to call
	Description string      `json:"description"`  // a brief description of what the tool does
	InputSchema InputSchema `json:"input_schema"` // defining what the tool uses as input
	/*
		Defines the function to be executed for the model's tool calls
	*/
	Function func(...any) any `json:"-"`
}

func (llm LLM) getToolOutput(response LLMResponse, tools []Tool) any {
	respContent := response.Content

	// Creating a map containing every tool's name for efficient lookup
	toolMap := make(map[string]Tool)
	for i := range tools {
		toolMap[tools[i].Name] = tools[i]
	}

	for j := range respContent {
		if respContent[j].Type == "tool_use" {
			tool, ok := toolMap[respContent[j].Name]
			if !ok {
				return nil
			}
			// Map that represents the name of each input parameter and the value to pass in for that parameter for this function call
			funcParams := respContent[j].Input

			// Creating a slice of function parameter values from the given input parameters map
			values := make([]any, 0, len(funcParams))
			for _, value := range funcParams {
				values = append(values, value)
			}
			return tool.Function(values...)

		}
	}
	return nil
}

/*
Tool for searching the web and retrieving relevant search results
*/
type SearchTool struct {
	apiKey      string // Brave API Key
	name        string
	description string
}

/*
Tool for querying a PostgreSQL database
*/
type SQLTool struct {
}
