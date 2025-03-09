package main

// TODO: replace w/ MongoDB?
type MemoryBlock struct {
	UserMemory/*map[string]string*/ string
	AgentMemory/*map[string]string*/ string
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
*/
func (memoryBlock MemoryBlock) coreMemorySave(agentMem bool, section string, memory string) MemoryBlock {
	if agentMem {
		memoryBlock.AgentMemory += "\n"
		memoryBlock.AgentMemory += memory
		return memoryBlock
	}
	memoryBlock.UserMemory += "\n"
	memoryBlock.UserMemory += memory
	return memoryBlock
}

// func (llm *LLM) addMemoryUpdateToChatHistory(content Content, memoryBlock MemoryBlock) {
// 	llm.Messages = append(llm.Messages,
// 		Message{Role: "user",
// 			Content: []Content{
// 				{
// 					Type:      "tool_result",
// 					ToolUseID: content.ToolUseID,
// 					Name:      content.Name,
// 					Content:   fmt.Sprintf("Updated Memory: %+v\n", memoryBlock),
// 				},
// 			},
// 		},
// 	)
// }
