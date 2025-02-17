package main

type Tool interface {
}

type Agent struct {
	model   string
	role    string
	tools   []Tool
	context string
}
