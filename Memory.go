package main

import (
	"strings"
)

// TODO: ability to add custom memory blocks into core memory (aside from just human and agent) and corresponding tools for those memory blocks
type MemoryBlock struct {
	sectionName string
	size        int           // number of characters allotted for a given block of memory
	data        string        // the text data stored in the given block of memory
	toString    func() string // converting the data in the memory block into a string that can be parsed by the Agent
}

// TODO: replace w/ ElasticSearch?
type CoreMemory struct {
	UserMemory/*map[string]string*/ MemoryBlock
	AgentMemory/*map[string]string*/ MemoryBlock
}

func NewCoreMemoryUnit() *CoreMemory {
	coreMemory := CoreMemory{
		UserMemory:  MemoryBlock{sectionName: "user", size: 2000, data: ""},
		AgentMemory: MemoryBlock{sectionName: "agent", size: 2000, data: ""},
	}

	// coreMemory.UserMemory.toString = func() string {

	// }

	// coreMemory.AgentMemory.toString = func() string {

	// }

	return &coreMemory
}

/*
Conversation history
*/
type RecallMemory struct {
}

/*
General purpose out of context information storage
Retrieve via RAG (semantic search)
*/
type ArchivalMemory struct {
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
		llm.CoreMemory.AgentMemory.data += "\n"
		llm.CoreMemory.AgentMemory.data += newContent
		return llm.CoreMemory
	}
	llm.CoreMemory.UserMemory.data += "\n"
	llm.CoreMemory.UserMemory.data += newContent
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

func (llm *Agent) coreMemoryReplace(agentMem bool, oldContent string, newContent string, requestHeartbeat bool) CoreMemory {
	llm.HeartbeatState = requestHeartbeat

	if agentMem {
		llm.CoreMemory.AgentMemory.data = strings.ReplaceAll(llm.CoreMemory.AgentMemory.data, oldContent, newContent)
		return llm.CoreMemory
	}
	llm.CoreMemory.UserMemory.data = strings.ReplaceAll(llm.CoreMemory.UserMemory.data, oldContent, newContent)
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

func (llm Agent) recursiveSummarization() {

}

func (llm *Agent) flushMemory() {

}

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
