package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cassamajor/agentic"
)

func main() {
	tools := []agentic.ToolDefinition{agentic.ReadFileDefinition, agentic.ListFilesDefinition, agentic.EditFileDefinition}

	agent, err := agentic.NewAgent(
		agentic.WithTools(tools),
	)

	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
	}
}
