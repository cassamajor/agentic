package main

type Role string

const (
	User         Role = "User"
	Assistant    Role = "Assistant"
	System       Role = "System"
	ToolCall     Role = "Tool Call"
	ToolResponse Role = "Tool Response"
)

type Message struct {
	Role    Role
	Content string
}

type Messages []Message
