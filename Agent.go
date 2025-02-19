package main

type Agent struct {
	llm    LLM    // the name of the model to use
	role   string // the role that this agent has in the workflow
	tools  []Tool // the list of tools that the agent has access to
	memory string // storing information that the agent is to use as part of its recursive calls
	cont   bool   // used to determine if the agent is to continue invoking itself or to stop
}
