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
TODO: should this return the entire core memory, or just the updated block of memory?
*/
func (llm *Agent) coreMemoryAppend(section string, newContent string, requestHeartbeat bool) *CoreMemory {
	llm.HeartbeatState = requestHeartbeat
	memoryBlock, ok := llm.CoreMemory.Blocks[section]
	if !ok {
		return nil
	}

	memoryBlock.data += fmt.Sprintf("\n%s", newContent)
	return &llm.CoreMemory
}

/*
Helper function defined to be used as a Tool for an agent
*/
func (llm *Agent) coreMemoryAppendVariadic(args ...any) any {
	if len(args) == 0 {
		return "Error: No parameters provided"
	}

	section, ok := args[0].(string)

	if !ok {
		return "Error: Invalid section string"
	}

	newContent, ok := args[1].(string)

	if !ok {
		return "Error: Invalid newContent string"
	}

	requestHeartbeat, ok := args[2].(bool)

	if !ok {
		return "Error: Invalid requestHeartbeat flag"
	}

	return llm.coreMemoryAppend(section, newContent, requestHeartbeat)
}

func (llm Agent) createCoreMemoryAppendTool() *Tool {
	return &Tool{
		Name:        "coreMemoryAppend",
		Description: "Save important information about you (the agent) or the human you are chatting with, inside of your core memory.",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"section": {
					Type:        "string",
					Description: "Represents the section of core memory to append to. This could be \"User\" (to store information about the user you are talking with), \"Agent\" (to store information about yourself, the agent), or some other user-defined portion of your core memory.",
				},
				"newContent": {
					Type:        "string",
					Description: "The text information to store in your core memory that you will be able to refer to in the future.",
				},
				"requestHeartbeat": {
					Type:        "boolean",
					Description: "Set to `true` to continue your execution in a loop (i.e. you will be invoked again), allowing you to think once again. Set to `false` to end your execution now.",
				},
			},
			Required: []string{"section", "newContent", "requestHeartbeat"},
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

/* Core Memory Replace */

func (llm *Agent) coreMemoryReplace(section string, oldContent string, newContent string, requestHeartbeat bool) *CoreMemory {
	llm.HeartbeatState = requestHeartbeat
	memoryBlock, ok := llm.CoreMemory.Blocks[section]
	if !ok {
		return nil
	}

	memoryBlock.data = strings.ReplaceAll(memoryBlock.data, oldContent, newContent)

	return &llm.CoreMemory
}

/*
Helper function defined to be used as a Tool for an agent
*/
func (llm *Agent) coreMemoryReplaceVariadic(args ...any) any {
	if len(args) == 0 {
		return "Error: No parameters provided"
	}

	section, ok := args[0].(string)

	if !ok {
		return "Error: Invalid section flag"
	}

	oldContent, ok := args[1].(string)

	if !ok {
		return "Error: Invalid oldContent string"
	}

	newContent, ok := args[2].(string)

	if !ok {
		return "Error: Invalid newContent string"
	}

	requestHeartbeat, ok := args[3].(bool)

	if !ok {
		return "Error: Invalid requestHeartbeat flag"
	}

	return llm.coreMemoryReplace(section, oldContent, newContent, requestHeartbeat)
}

func (llm Agent) createCoreMemoryReplaceTool() *Tool {
	return &Tool{
		Name:        "coreMemoryReplace",
		Description: "Takes a piece of existing information about you (the agent) or the human you are chatting with, and replaces this information with a new piece of information.",
		//Description: "Save important information about you (the agent) or the human you are chatting with.",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"section": {
					Type:        "string",
					Description: "The section of core memory to update. This could be \"User\" (to update information about the user you are talking with), \"Agent\" (to update information about yourself, the agent), or some other user-defined portion of your core memory.",
				},
				"oldContent": {
					Type:        "string",
					Description: "The text information in your core memory that is now outdated and is to be replaced.",
				},
				"newContent": {
					Type:        "string",
					Description: "The new text information to store in your core memory to replace an existing piece of information.",
				},
				"requestHeartbeat": {
					Type:        "boolean",
					Description: "Set to `true` to continue your execution in a loop (i.e. you will be invoked again), allowing you to think once again. Set to `false` to end your execution now.",
				},
			},
			Required: []string{"section", "oldContent", "newContent", "requestHeartbeat"},
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
func (llm *Agent) pauseHeartbeats() bool {
	llm.HeartbeatState = false
	return true
}

func (llm *Agent) createPauseHeartbeatsTool() *Tool {
	return &Tool{
		Name:        "pauseHeartbeats",
		Description: "Pauses your execution.",
		//Description: "Save important information about you (the agent) or the human you are chatting with.",
		InputSchema: InputSchema{
			Type:       "object",
			Properties: map[string]Property{},
			Required:   []string{},
		},
		Function: func(args ...any) any {
			return llm.pauseHeartbeats()
		},
	}
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
