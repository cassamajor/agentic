package agentic

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

type BashInput struct {
	Command string `json:"cmd" jsonschema:"Command (including arguments) to execute in the terminal."`
	DryRun  bool   `json:"dryRun" jsonschema:"If set to true, the bash tool will returns the command string without executing it."`
}

var BashInputSchema = GenerateSchema[BashInput]()

func BashTool(input json.RawMessage) (string, error) {
	bashInput := BashInput{}
	err := json.Unmarshal(input, &bashInput)

	if err != nil {
		return "", err
	}
	cmd, err := CmdFromString(bashInput.Command)

	if bashInput.DryRun {
		return cmd.String(), nil
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

var BashToolDefinition = ToolDefinition{
	Name: "bash",
	Description: `Execute a command into the operating system terminal.
	When DryRun is specified, the bash tool returns the command string without executing it.`,
	InputSchema: BashInputSchema,
	Function:    BashTool,
}

func CmdFromString(input string) (*exec.Cmd, error) {
	args := strings.Fields(input)

	// error if args is empty
	if len(args) == 0 {
		return nil, errors.New("input cannot be empty")
	}

	return exec.Command(args[0], args[1:]...), nil
}
