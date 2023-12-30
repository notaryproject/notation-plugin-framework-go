package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unsafe"

	"github.com/notaryproject/notation-plugin-framework-go/internal/slices"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

func getValidArgsString(md *plugin.GetMetadataResponse) string {
	return fmt.Sprintf(`<%s>`, strings.Join(getValidArgs(md), "|"))
}

// getValidArgs returns list of valid arguments depending upon the plugin capabilities
func getValidArgs(md *plugin.GetMetadataResponse) []string {
	opts := []plugin.Command{
		plugin.CommandGetMetadata, plugin.Version,
	}

	if slices.Contains(md.Capabilities, plugin.CapabilitySignatureGenerator) {
		opts = append(opts, plugin.CommandGenerateSignature, plugin.CommandDescribeKey)
	}

	if slices.Contains(md.Capabilities, plugin.CapabilityEnvelopeGenerator) {
		opts = append(opts, plugin.CommandGenerateEnvelope)
	}

	if slices.Contains(md.Capabilities, plugin.CapabilityTrustedIdentityVerifier) || slices.Contains(md.Capabilities, plugin.CapabilityRevocationCheckVerifier) {
		opts = append(opts, plugin.CommandVerifySignature)
	}
	// convert &value (type *[]Command) to *[]string via unsafe.Pointer, then deref
	stringOpts := *(*[]string)(unsafe.Pointer(&opts))
	sort.Strings(stringOpts)
	return stringOpts
}

// deliverError print to standard error and then return nonzero exit code
func deliverError(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
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
