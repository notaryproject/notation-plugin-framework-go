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
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/notaryproject/notation-plugin-framework-go/internal/slices"
	"github.com/notaryproject/notation-plugin-framework-go/log"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

// CLI struct is used to create an executable for plugin.
type CLI struct {
	name   string
	pl     plugin.Plugin
	logger log.Logger
}

// New creates a new CLI using given plugin
func New(pl plugin.Plugin) (*CLI, error) {
	return NewWithLogger(pl, &discardLogger{})
}

// NewWithLogger creates a new CLI using given plugin and logger
func NewWithLogger(pl plugin.Plugin, l log.Logger) (*CLI, error) {
	if pl == nil {
		return nil, errors.New("plugin cannot be nil")
	}

	return &CLI{
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
		deliverError(fmt.Sprintf("Invalid command, valid commands are: %s", getValidArgsString(md)))
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
