package main

import (
	"fmt"
	"strings"
)

// TODO: ability to add custom memory blocks into core memory (aside from just human and agent) and corresponding tools for those memory blocks
type MemoryBlock struct {
	sectionName string
	//size        int           // number of characters allotted for a given block of memory
	data string // the text data stored in the given block of memory
}

// TODO: replace w/ ElasticSearch?
type CoreMemory struct {
	Blocks map[string]*MemoryBlock
	// UserMemory/*map[string]string*/ MemoryBlock
	// AgentMemory/*map[string]string*/ MemoryBlock
}

func NewCoreMemoryUnit(persona string) *CoreMemory {
	coreMemory := CoreMemory{Blocks: map[string]*MemoryBlock{
		"User":  {sectionName: "user" /*size: 2000,*/, data: ""},
		"Agent": {sectionName: "agent" /*size: 2000,*/, data: persona}},
	}

	// coreMemory.UserMemory.toString = func() string {

	// }

	// coreMemory.AgentMemory.toString = func() string {

	// }

	return &coreMemory
}

/*
Conversation history -- stored out of context
*/
type RecallMemory struct {
}

/*
General purpose out of context information storage
Retrieve via RAG (semantic search/tf-idf)
*/
type ArchivalMemory struct {
}

/*
Converts the data stored across all of Core Memory into a string that can be used in the Agent's context window.
*/
func (coreMemory CoreMemory) toString() string {
	var builder strings.Builder

	for section, block := range coreMemory.Blocks {
		builder.WriteString(fmt.Sprintf("%s:\n%s\n", section, block.data))
	}
	return builder.String()
}

// func (memoryBlock *MemoryBlock) coreMemorySave(agentMem bool, section string, memory string) (result struct {
// 	sec string
// 	mem string
// }) {
// 	if agentMem {
// 		memoryBlock.AgentMemory += "\n"
// 		memoryBlock.AgentMemory += memory
// 		result.sec = section
// 		result.mem = memory
// 		return
// 	}
// 	memoryBlock.UserMemory += "\n"
// 	memoryBlock.UserMemory += memory
// 	result.sec = section
// 	result.mem = memory
// 	return
// }

/*
Stores a memory (either User or Agent memory) in the memory block.
Ex. memoryBlock = memoryBlock.coreMemorySave(false, "name", "Bob")
TODO: should this return the entire memory block, or just the updated portion of memory?
*/
func (llm *Agent) coreMemoryAppend(agentMem bool, newContent string, requestHeartbeat bool) CoreMemory {
	llm.HeartbeatState = requestHeartbeat

	if agentMem {
		llm.CoreMemory.Blocks["Agent"].data += fmt.Sprintf("\n%s", newContent)
		return llm.CoreMemory
	}
	llm.CoreMemory.Blocks["User"].data += fmt.Sprintf("\n%s", newContent)
	return llm.CoreMemory
}

/*
Helper function defined to be used as a Tool for an agent
*/
func (llm *Agent) coreMemoryAppendVariadic(args ...any) any {
	if len(args) == 0 {
		return "Error: No parameters provided"
	}

	agentMem, ok := args[0].(bool)

	if !ok {
		return "Error: Invalid section flag"
	}

	newContent, ok := args[1].(string)

	if !ok {
		return "Error: Invalid newContent string"
	}

	requestHeartbeat, ok := args[2].(bool)

	if !ok {
		return "Error: Invalid requestHeartbeat flag"
	}

	return llm.coreMemoryAppend(agentMem, newContent, requestHeartbeat)
}

func (llm Agent) createCoreMemoryAppendTool() *Tool {
	return &Tool{
		Name:        "coreMemoryAppend",
		Description: "Save important information about you (the agent) or the human you are chatting with.",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"agentMem": {
					Type:        "boolean",
					Description: "Must be either `true` to save information about yourself (the agent), or `false` to save information about the user.",
				},
				"newContent": {
					Type:        "string",
					Description: "CoreMemory to save in the section",
				},
				"requestHeartbeat": {
					Type:        "boolean",
					Description: "Set to `true` to continue your execution (i.e. you will be invoked again) and run an additional step afterward, or `false` to end your execution now.",
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
	}
}

func (llm *Agent) coreMemoryReplace(agentMem bool, oldContent string, newContent string, requestHeartbeat bool) CoreMemory {
	llm.HeartbeatState = requestHeartbeat

	if agentMem {
		llm.CoreMemory.Blocks["Agent"].data = strings.ReplaceAll(llm.CoreMemory.Blocks["Agent"].data, oldContent, newContent)
		return llm.CoreMemory
	}
	llm.CoreMemory.Blocks["User"].data = strings.ReplaceAll(llm.CoreMemory.Blocks["User"].data, oldContent, newContent)
	return llm.CoreMemory
}

func (llm *Agent) conversationSearch(requestHeartbeat bool) {
	llm.HeartbeatState = requestHeartbeat

}

func (llm *Agent) archivalSearch(query string, requestHeartbeat bool) /*[]string*/ {
	llm.HeartbeatState = requestHeartbeat

}

func (llm *Agent) archivalInsert(newContent string, requestHeartbeat bool) {
	llm.HeartbeatState = requestHeartbeat

}

/*
Pauses the agent loop
*/
func (llm *Agent) pauseHeartbeats() {
	llm.HeartbeatState = false
}

/* Functions to clear space in the context window */

// func (llm Agent) recursiveSummarization() {

// }

// func (llm *Agent) flushMemory() {

// }

// func (llm *Agent) addMemoryUpdateToChatHistory(content Content, memoryBlock MemoryBlock) {
// 	llm.Messages = append(llm.Messages,
// 		Message{Role: "user",
// 			Content: []Content{
// 				{
// 					Type:      "tool_result",
// 					ToolUseID: content.ToolUseID,
// 					Name:      content.Name,
// 					Content:   fmt.Sprintf("Updated CoreMemory: %+v\n", memoryBlock),
// 				},
// 			},
// 		},
// 	)
// }
