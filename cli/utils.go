package cli

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/notaryproject/notation-plugin-framework-go/internal/slices"
	"github.com/notaryproject/notation-plugin-framework-go/log"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

func getValidArgsString(md *plugin.GetMetadataResponse) string {
	return fmt.Sprintf(`<%s>`, strings.Join(getValidArgs(md), "|"))
}

// getValidArgs returns list of valid arguments depending upon the plugin capabilities
func getValidArgs(md *plugin.GetMetadataResponse) []string {
	opts := []string{
		"get-plugin-metadata", "version",
	}

	if slices.Contains(md.Capabilities, plugin.CapabilitySignatureGenerator) {
		opts = append(opts, "generate-signature", "describe-key")
	}

	if slices.Contains(md.Capabilities, plugin.CapabilityEnvelopeGenerator) {
		opts = append(opts, "generate-envelope")
	}

	if slices.Contains(md.Capabilities, plugin.CapabilityTrustedIdentityVerifier) || slices.Contains(md.Capabilities, plugin.CapabilityRevocationCheckVerifier) {
		opts = append(opts, "verify-signature")
	}
	sort.Strings(opts)
	return opts
}

func getMetadata(ctx context.Context, p plugin.Plugin) *plugin.GetMetadataResponse {
	md, err := p.GetMetadata(ctx, &plugin.GetMetadataRequest{})
	if err != nil {
		logger := log.GetLogger(ctx)
		logger.Errorf("GetMetadataRequest error :%v", err)
		deliverError("Error: Failed to get plugin metadata.")
	}
	return md
}

// deliverError print to standard error and then return nonzero exit code
func deliverError(message string) {
	_, _ = fmt.Fprintf(os.Stderr, message)
	os.Exit(1)
}
