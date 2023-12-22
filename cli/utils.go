package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/notaryproject/notation-plugin-framework-go/internal/slices"
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

// deliverError print to standard error and then return nonzero exit code
func deliverError(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
