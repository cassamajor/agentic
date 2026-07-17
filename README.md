Batteries included agent. In your shell, run:

```shell
go get github.com/cassamajor/agentic
```

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cassamajor/agentic"
)

func main() {
	agent, err := agentic.NewAgent()

	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
	}
}
```

Based the Amp article on [How to Build an Agent](https://ampcode.com/notes/how-to-build-an-agent), with additional optimizations.

### Thoughts & TODOs
The `mcp.Tool` struct only requires `Name` and `Description`. It has a default InputSchema of `map[string]any`.
`mcp.AddTool` requires the `mcp.Tool` struct and the tool function.

An MCP tool function must satisfy the signature: `context.Context, *CallToolRequest, any`.

Our custom `ToolDefinition` also has `Name`, `Description`, and `InputField`. Additionally, it also takes the `Function` directly.

It would be trivial to map our custom `ToolDefinition` to `mcp.AddTool`. We would need to change the shape of our tool function, however.

Each tool function currently has the signature `input json.RawMessage` and returns `(string, error)`. It would need to instead match `context.Context, *CallToolRequest, any` and return `*CallToolResult, any, error`.

I can write an adapter for this. `executeTool` is impacted by this change and would need to be updated.

These changes should be handled within `WithTools`.

I should move the tooling logic from `runInference` to `WithTools`. The type for `Agent.Tools` should be updated to `anthropic.ToolUnionParam` so the converted tools can persist.
At that point I can actually get rid of the `runInference` function and move the logic into the `Run` function.

MCP's equivalent of `GenerateSchema` is `setSchema`.

This package needs to be updated to use `github.com/google/jsonschema-go/jsonschema` rather than `github.com/invopop/jsonschema`.
`reflect` is part of the standard library.