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