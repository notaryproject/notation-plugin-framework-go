// Package cli provides boilerplate code required to generate plugin executable.
// At high level it performs following tasks
// 1. Validate command arguments
// 2. Read and unmarshal input
// 3. Execute relevant plugin functions
// 4. marshals output

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/notaryproject/notation-plugin-framework-go/internal/slices"
	"github.com/notaryproject/notation-plugin-framework-go/log"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

type CLI struct {
	name string
	pl   plugin.Plugin
}

// New creates a new CLI using given plugin
func New(executableName string, pl plugin.Plugin) *CLI {
	return &CLI{
		name: executableName,
		pl:   pl}
}

// Execute is main controller that reads/validates commands, parses input, executes relevant plugin functions
// and returns corresponding output.
func (c CLI) Execute(ctx context.Context, args []string) {
	c.validateArgs(ctx, args)

	rescueStdOut := deferStdout()
	command := args[1]
	var resp any
	var err error
	logger := log.GetLogger(ctx)
	switch plugin.Command(command) {
	case plugin.CommandGetMetadata:
		var request plugin.GetMetadataRequest
		err = unmarshalRequest(ctx, &request)
		if err == nil {
			logger.Debugf("executing %s plugin's GetMetadata function", reflect.TypeOf(c.pl))
			resp, err = c.pl.GetMetadata(ctx, &request)
		}
	case plugin.CommandGenerateEnvelope:
		var request plugin.GenerateEnvelopeRequest
		err = unmarshalRequest(ctx, &request)
		if err == nil {
			logger.Debugf("executing %s plugin's GenerateEnvelope function", reflect.TypeOf(c.pl))
			resp, err = c.pl.GenerateEnvelope(ctx, &request)
		}
	case plugin.CommandVerifySignature:
		var request plugin.VerifySignatureRequest
		err = unmarshalRequest(ctx, &request)
		if err == nil {
			logger.Debugf("executing %s plugin's VerifySignature function", reflect.TypeOf(c.pl))
			resp, err = c.pl.VerifySignature(ctx, &request)
		}
	case plugin.CommandDescribeKey:
		var request plugin.DescribeKeyRequest
		err = unmarshalRequest(ctx, &request)
		if err == nil {
			logger.Debugf("executing %s plugin's DescribeKey function", reflect.TypeOf(c.pl))
			resp, err = c.pl.DescribeKey(ctx, &request)
		}
	case plugin.CommandGenerateSignature:
		var request plugin.VerifySignatureRequest
		err = unmarshalRequest(ctx, &request)
		if err == nil {
			logger.Debugf("executing %s plugin's VerifySignature function", reflect.TypeOf(c.pl))
			resp, err = c.pl.VerifySignature(ctx, &request)
		}
	case plugin.Version:
		rescueStdOut()
		c.printVersion(ctx)
	default:
		// should never happen
		rescueStdOut()
		deliverError(plugin.NewGenericError("something went wrong").Error())
	}

	op, pluginErr := marshalResponse(ctx, resp, err)
	rescueStdOut()
	if pluginErr != nil {
		deliverError(pluginErr.Error())
	}
	fmt.Println(op)
}

// printVersion prints version of executable
func (c CLI) printVersion(ctx context.Context) {
	md := getMetadata(ctx, c.pl)

	fmt.Printf("%s - %s\nVersion: %s", md.Name, md.Description, md.Version)
	os.Exit(0)
}

// validateArgs validate commands/arguments passed to executable.
func (c CLI) validateArgs(ctx context.Context, args []string) {
	md := getMetadata(ctx, c.pl)
	if !(len(args) == 2 && slices.Contains(getValidArgs(md), args[1])) {
		deliverError(fmt.Sprintf("Invalid command, valid choices are: %s %s", c.name, getValidArgsString(md)))
	}
}

// deferStdout is used to make sure that nothing get emitted to stdout and stderr until intentionally rescued.
// This is required to make sure that the plugin or its dependency doesn't interferes with notation <-> plugin communication
func deferStdout() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null

	return func() {
		err := null.Close()
		if err != nil {
			return
		}
		os.Stdout = sout
		os.Stderr = serr
	}
}

// unmarshalRequest reads input from std.in and unmarshal it into given request struct
func unmarshalRequest(ctx context.Context, request any) error {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		logger := log.GetLogger(ctx)
		logger.Errorf("%s unmarshalling error :%v", reflect.TypeOf(request), err)
		return plugin.NewJSONParsingError(plugin.ErrorMsgMalformedInput)
	}
	return nil
}

// marshalResponse marshals the given response struct into json
func marshalResponse(ctx context.Context, response any, err error) (string, *plugin.Error) {
	logger := log.GetLogger(ctx)
	if err != nil {
		logger.Errorf("%s error :%v", reflect.TypeOf(response), err)
		if plgErr, ok := err.(*plugin.Error); ok {
			return "", plgErr
		}
		return "", plugin.NewGenericError(err.Error())
	}

	logger.Debug("marshalling response")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("%s marshalling error: %v", reflect.TypeOf(response), err)
		return "", plugin.NewGenericErrorf(plugin.ErrorMsgMalformedOutputFmt, err.Error())
	}

	logger.Debugf("%s response :%s", reflect.TypeOf(response), jsonResponse)
	return string(jsonResponse), nil
}
