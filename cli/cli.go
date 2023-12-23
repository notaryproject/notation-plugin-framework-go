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
	"strings"

	"github.com/notaryproject/notation-plugin-framework-go/internal/slices"
	"github.com/notaryproject/notation-plugin-framework-go/log"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

type CLI struct {
	name   string
	pl     plugin.Plugin
	logger log.Logger
}

// New creates a new CLI using given plugin
func New(executableName string, pl plugin.Plugin) (*CLI, error) {
	return NewWithLogger(executableName, pl, &discardLogger{})
}

// NewWithLogger creates a new CLI using given plugin and logger
func NewWithLogger(executableName string, pl plugin.Plugin, l log.Logger) (*CLI, error) {
	if strings.HasPrefix(executableName, plugin.BinaryPrefix) {
		return nil, fmt.Errorf("executable name must start with prefix: %s", plugin.BinaryPrefix)
	}

	return &CLI{
		name:   executableName,
		pl:     pl,
		logger: l,
	}, nil
}

// Execute is main controller that reads/validates commands, parses input, executes relevant plugin functions
// and returns corresponding output.
func (c *CLI) Execute(ctx context.Context, args []string) {
	c.validateArgs(ctx, args)

	rescueStdOut := deferStdout()
	command := args[1]
	var resp any
	var err error
	switch plugin.Command(command) {
	case plugin.CommandGetMetadata:
		var request plugin.GetMetadataRequest
		err = c.unmarshalRequest(&request)
		if err == nil {
			c.logger.Debugf("executing %s plugin's GetMetadata function", reflect.TypeOf(c.pl))
			resp, err = c.pl.GetMetadata(ctx, &request)
		}
	case plugin.CommandGenerateEnvelope:
		var request plugin.GenerateEnvelopeRequest
		err = c.unmarshalRequest(&request)
		if err == nil {
			c.logger.Debugf("executing %s plugin's GenerateEnvelope function", reflect.TypeOf(c.pl))
			resp, err = c.pl.GenerateEnvelope(ctx, &request)
		}
	case plugin.CommandVerifySignature:
		var request plugin.VerifySignatureRequest
		err = c.unmarshalRequest(&request)
		if err == nil {
			c.logger.Debugf("executing %s plugin's VerifySignature function", reflect.TypeOf(c.pl))
			resp, err = c.pl.VerifySignature(ctx, &request)
		}
	case plugin.CommandDescribeKey:
		var request plugin.DescribeKeyRequest
		err = c.unmarshalRequest(&request)
		if err == nil {
			c.logger.Debugf("executing %s plugin's DescribeKey function", reflect.TypeOf(c.pl))
			resp, err = c.pl.DescribeKey(ctx, &request)
		}
	case plugin.CommandGenerateSignature:
		var request plugin.VerifySignatureRequest
		err = c.unmarshalRequest(&request)
		if err == nil {
			c.logger.Debugf("executing %s plugin's VerifySignature function", reflect.TypeOf(c.pl))
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

	op, pluginErr := c.marshalResponse(resp, err)
	rescueStdOut()
	if pluginErr != nil {
		deliverError(pluginErr.Error())
	}
	fmt.Println(op)
}

// printVersion prints version of executable
func (c *CLI) printVersion(ctx context.Context) {
	md := c.getMetadata(ctx, c.pl)

	fmt.Printf("%s - %s\nVersion: %s \n", md.Name, md.Description, md.Version)
	os.Exit(0)
}

// validateArgs validate commands/arguments passed to executable.
func (c *CLI) validateArgs(ctx context.Context, args []string) {
	md := c.getMetadata(ctx, c.pl)
	if !(len(args) == 2 && slices.Contains(getValidArgs(md), args[1])) {
		deliverError(fmt.Sprintf("Invalid command, valid choices are: %s %s", c.name, getValidArgsString(md)))
	}
}

// unmarshalRequest reads input from std.in and unmarshal it into given request struct
func (c *CLI) unmarshalRequest(request any) error {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		c.logger.Errorf("%s unmarshalling error :%v", reflect.TypeOf(request), err)
		return plugin.NewJSONParsingError(plugin.ErrorMsgMalformedInput)
	}
	return nil
}

func (c *CLI) getMetadata(ctx context.Context, p plugin.Plugin) *plugin.GetMetadataResponse {
	md, err := p.GetMetadata(ctx, &plugin.GetMetadataRequest{})
	if err != nil {
		c.logger.Errorf("GetMetadataRequest error :%v", err)
		deliverError("Error: Failed to get plugin metadata.")
	}
	return md
}

// marshalResponse marshals the given response struct into json
func (c *CLI) marshalResponse(response any, err error) (string, *plugin.Error) {
	if err != nil {
		c.logger.Errorf("%s error: %v", reflect.TypeOf(response), err)
		if plgErr, ok := err.(*plugin.Error); ok {
			return "", plgErr
		}
		return "", plugin.NewGenericError(err.Error())
	}

	c.logger.Debug("marshalling response")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		c.logger.Errorf("%s marshalling error: %v", reflect.TypeOf(response), err)
		return "", plugin.NewGenericErrorf(plugin.ErrorMsgMalformedOutputFmt, err.Error())
	}

	c.logger.Debugf("%s response: %s", reflect.TypeOf(response), jsonResponse)
	return string(jsonResponse), nil
}

// deferStdout is used to make sure that nothing get emitted to stdout and stderr until intentionally rescued.
// This is required to make sure that the plugin or its dependency doesn't interfere with notation <-> plugin communication
func deferStdout() func() {
	// Ignoring error because we don't want plugin to fail if `os.DevNull` is misconfigured.
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

// discardLogger implements Logger but logs nothing. It is used when user
// disenabled logging option in notation, i.e. loggerKey is not in the context.
type discardLogger struct{}

func (dl *discardLogger) Debug(_ ...interface{}) {
}

func (dl *discardLogger) Debugf(_ string, _ ...interface{}) {
}

func (dl *discardLogger) Debugln(_ ...interface{}) {
}

func (dl *discardLogger) Info(_ ...interface{}) {
}

func (dl *discardLogger) Infof(_ string, _ ...interface{}) {
}

func (dl *discardLogger) Infoln(_ ...interface{}) {
}

func (dl *discardLogger) Warn(_ ...interface{}) {
}

func (dl *discardLogger) Warnf(_ string, _ ...interface{}) {
}

func (dl *discardLogger) Warnln(_ ...interface{}) {
}

func (dl *discardLogger) Error(_ ...interface{}) {
}

func (dl *discardLogger) Errorf(_ string, _ ...interface{}) {
}

func (dl *discardLogger) Errorln(_ ...interface{}) {
}
