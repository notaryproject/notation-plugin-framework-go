package main

import (
	"context"
	"fmt"
	"os"

	"github.com/notaryproject/notation-plugin-framework-go/cli"
)

func main() {
	ctx := context.Background()
	// Initialize plugin
	plugin, err := NewExamplePlugin()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize plugin")
		os.Exit(2)
	}

	// Create executable
	cli.New(plugin).Execute(ctx, os.Args)
}
