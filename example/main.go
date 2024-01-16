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
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize plugin: %v\n", err)
		os.Exit(2)
	}

	// Create executable
	pluginCli, err := cli.New(plugin)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create executable: %v\n", err)
		os.Exit(3)
	}
	pluginCli.Execute(ctx, os.Args)
}
