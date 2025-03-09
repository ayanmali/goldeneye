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
TODO: should this return the entire memory block, or just the updated portion of memory?
*/
func (memoryBlock MemoryBlock) coreMemorySave(agentMem bool, memory string) MemoryBlock {
	if agentMem {
		memoryBlock.AgentMemory += "\n"
		memoryBlock.AgentMemory += memory
		return memoryBlock
	}
	memoryBlock.UserMemory += "\n"
	memoryBlock.UserMemory += memory
	return memoryBlock
}

/*
Helper function defined to be used as a Tool for an agent
*/
func (memoryBlock MemoryBlock) coreMemorySaveVariadic(args ...any) any {
	if len(args) == 0 {
		return "Error: No parameters provided"
	}

	agentMem, ok := args[0].(bool)

	if !ok {
		return "Error: Invalid section flag"
	}

	memory, ok := args[1].(string)

	if !ok {
		return "Error: Invalid memory string"
	}

	return memoryBlock.coreMemorySave(agentMem, memory)
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
