package agentic

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

const (
	USER string = "\x1b[94mYou\x1b[0m: "          // Colorize the text "You"
	LLM  string = "\x1b[93mLLM\x1b[0m: %s\n"      // Colorize the text "LLM"
	TOOL string = "\x1b[92mtool\x1b[0m: %s(%s)\n" // Colorize the text "tool"
)

// Description should follow best practices: a brief explanation, specifiy the circumstances the tool should be used, and the circumstances that it should not be used.
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

type Agent struct {
	client       *anthropic.Client
	model        string
	SystemPrompt string
	UserInput    io.Reader
	Output       io.Writer
	Tools        []ToolDefinition
}

func (a *Agent) executeTool(content *anthropic.ContentBlockUnion) anthropic.ContentBlockParamUnion {
	var toolDef ToolDefinition
	var found bool

	for _, tool := range a.Tools {
		if tool.Name == content.Name {
			toolDef = tool
			found = true
			break
		}
	}
	if !found {
		return anthropic.NewToolResultBlock(content.ID, "tool not found", true)
	}

	fmt.Fprintf(a.Output, TOOL, content.Name, content.Input)

	// Call the function assigned to the tool definition.
	response, err := toolDef.Function(content.Input)

	// The tool function returned an error
	if err != nil {
		return anthropic.NewToolResultBlock(content.ID, err.Error(), true)
	}

	// Return the content produced by the tool function
	return anthropic.NewToolResultBlock(content.ID, response, false)
}

// runInference sends messages to the LLM and returns the response. It also specifies which tools are available to the agent.
func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	anthropicTools := []anthropic.ToolUnionParam{}

	for _, t := range a.Tools {
		oftool := &anthropic.ToolParam{
			Name:        t.Name,
			Description: anthropic.String(t.Description),
			InputSchema: t.InputSchema,
		}

		tool := anthropic.ToolUnionParam{OfTool: oftool}

		anthropicTools = append(anthropicTools, tool)
	}

	messageParams := anthropic.MessageNewParams{
		Model:     a.model,
		System:    []anthropic.TextBlockParam{{Text: a.SystemPrompt}},
		MaxTokens: int64(1042),
		Messages:  conversation,
		Tools:     anthropicTools,
	}

	message, err := a.client.Messages.New(ctx, messageParams)
	return message, err
}

// Run communicates with the LLM using an inner and outer loop.
// The outer loop handles the user's request, while the inner loop returns the LLM's response.
//
// Outer Loop:
//  1. Take input from the user and add it to the conversation slice.
//
// Inner Loop:
//  2. Send the conversation to the LLM.
//  3. Add the LLM's response to the conversation slice.
//  4. Print the LLM's response to the screen.
//  5. If there is a tool request, execute the tool and collect the response. Append the response to the conversation and resume from the beginning of the inner loop,
//     which will send the tool response to the LLM. The LLM can then react to the tool response without interaction from the user.
//  6. If there is not a tool request, then exit the inner loop and return to the outer loop.
func (a *Agent) Run(ctx context.Context) error {
	fmt.Fprintln(a.Output, "Chat with an LLM (use 'ctrl-c' to quit)")

	conversation := []anthropic.MessageParam{}

	// Collect user input
	scanner := bufio.NewScanner(a.UserInput)

	// Begin outer loop
	for {
		fmt.Fprint(a.Output, USER)

		if !scanner.Scan() {
			break // If there's no user input, exit the outer loop (exit the function call).
		}

		// Store user input
		userInput := scanner.Text()
		userMessage := anthropic.NewUserMessage(
			anthropic.NewTextBlock(userInput),
		)

		conversation = append(conversation, userMessage)

		// Begin inner loop
		for {
			// Send user input to Anthropic API and receive a response
			message, err := a.runInference(ctx, conversation)
			if err != nil {
				return err
			}

			// Add agent response to the conversation
			conversation = append(conversation, message.ToParam())

			toolResults := []anthropic.ContentBlockParamUnion{}

			// Print the agent response to the user
			for _, content := range message.Content {
				switch content.Type {
				case "text":
					fmt.Fprintf(a.Output, LLM, content.Text)
				case "tool_use":
					result := a.executeTool(&content)
					toolResults = append(toolResults, result)
				}
			}
			if len(toolResults) == 0 {
				break // exit the inner loop, return back to the outer loop
			}
			conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
		}
	}

	return nil
}

func NewAgent(opts ...option) (*Agent, error) {
	client := anthropic.NewClient()

	a := &Agent{
		client:    &client,
		UserInput: os.Stdin,
		Output:    os.Stdout,
		model:     anthropic.ModelClaudeOpus4_8,
	}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}
