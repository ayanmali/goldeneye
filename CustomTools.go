package main

/*
Tool for searching the web and retrieving relevant search results
*/
type WebSearchTool struct {
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
Tool for querying a MongoDB database
*/
type MongoDBTool struct {
}

/*
Tool for executing Go code
*/
type CodeExecutorTool struct {
}
