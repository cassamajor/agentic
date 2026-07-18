package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/cassamajor/agentic"
)

func main() {
	tools := []agentic.ToolDefinition{agentic.ReadFileDefinition, agentic.ListFilesDefinition, agentic.EditFileDefinition}

	client := anthropic.NewClient(
		option.WithBaseURL("http://localhost:8000"),
	)

	agent, err := agentic.NewAgent(
		agentic.WithTools(tools),
		agentic.WithClient(&client),
		agentic.WithModel("LFM2.5-8B-A1B-MLX-8bit"),
	)

	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
	}
}
