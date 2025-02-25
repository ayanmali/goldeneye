//
/* TOOL USE LOGIC */
//
package main

type Property struct {
	Type        string   `json:"type"`           // the data type of this property (parameter)
	Description string   `json:"description"`    // a brief description of what this property (parameter) represents
	Enum        []string `json:"enum:omitempty"` // possible values for Enum type input parameters
}

type InputSchema struct {
	// usually "object"
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

type Tool struct {
	Name string `json:"name"` // i.e. name of the function/API endpoint to call
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
	Description string      `json:"description"`
	InputSchema InputSchema `json:"input_schema"` // defining what the tool uses as input
	/*
		Defines the function to be executed for the model's tool calls
	*/
	Function func(...any) any `json:"-"`
}

// func (response LLMResponse) getToolOutput(request LLMRequest) any {
// 	respContent := response.Content
// 	// Getting the slice of Tools provided to the LLM in the initial request
// 	tools := request.Tools

// 	// Creating a map containing every tool's name for efficient lookup
// 	toolMap := make(map[string]Tool)
// 	for i := range tools {
// 		toolMap[tools[i].Name] = tools[i]
// 	}

// 	for j := range respContent {
// 		if respContent[j].Type == "tool_use" {
// 			tool, ok := toolMap[respContent[j].Name]
// 			if !ok {
// 				return nil
// 			}
// 			// Map that represents the name of each input parameter and the value to pass in for that parameter for this function call
// 			funcParams := respContent[j].Input

// 			// Creating a slice of function parameter values from the given input parameters map
// 			values := make([]any, 0, len(funcParams))
// 			for _, value := range funcParams {
// 				values = append(values, value)
// 			}
// 			return tool.Function(values...)

// 		}
// 	}
// 	return nil
// }

/*
Gets the result (output) of a single function call as specified by the LLM
*/
func getSingleToolOutput(content Content, request LLMRequest) (any, bool) {
	if content.Type != "tool_use" {
		return nil, false
	}

	// Getting the slice of Tools provided to the LLM in the initial request
	tools := request.Tools
	// Creating a map containing every tool's name for efficient lookup
	toolMap := make(map[string]Tool)
	for _, tool := range tools {
		toolMap[tool.Name] = tool
	}

	tool, ok := toolMap[content.Name]
	if !ok {
		return nil, false
	}
	// Map that represents the name of each input parameter and the value to pass in for that parameter for this function call
	funcParams := content.Input

	// Creating a slice of function parameter values from the given input parameters map
	values := make([]any, 0, len(funcParams))
	for _, value := range funcParams {
		values = append(values, value)
	}
	return tool.Function(values...), true
}

/*
Adds the result of executing a given tool to the LLM's context.
*/
func (request *LLMRequest) addToolResultToChatHistory(content Content) {
	messageToAppend := Message{Role: "user", Content: []Content{}}

	// Retrieving the output from the function call
	output, ok := getSingleToolOutput(content, *request) // perform type assertion
	if !ok {
		// Handle the case where the assertion fails
		output = "Error: Unable to retrieve tool output"
	}
	strOutput := output.(string)

	// Creating the message
	messageToAppend.Content = append(messageToAppend.Content, Content{
		Type:      "tool_result",
		ToolUseID: content.ToolUseID,
		Content:   strOutput,
	})

	// Adding the message to the request
	request.Messages = append(request.Messages, messageToAppend)
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

/*
Tool for executing Go code
*/
type CodeExecutorTool struct {
}
