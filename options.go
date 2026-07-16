package agentic

import (
	"errors"
	"io"

	"github.com/anthropics/anthropic-sdk-go"
)

type option func(*Agent) error

func WithClient(c *anthropic.Client) option {
	return func(a *Agent) error {
		if c == nil {
			return errors.New("nil is not a valid client")
		}
		a.client = c
		return nil
	}
}

func WithInput(r io.Reader) option {
	return func(a *Agent) error {
		if r == nil {
			return errors.New("nil is not a valid reader")
		}
		a.UserInput = r
		return nil
	}
}

func WithOutput(w io.Writer) option {
	return func(a *Agent) error {
		if w == nil {
			return errors.New("nil is not a valid writer")
		}
		a.Output = w
		return nil
	}
}

func WithTools(td []ToolDefinition) option {
	return func(a *Agent) error {
		if td == nil {
			return errors.New("nil is not a valid tool definition")
		}
		a.Tools = td
		return nil
	}
}
