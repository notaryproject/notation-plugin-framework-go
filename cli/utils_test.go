// Copyright The Notary Project Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"reflect"
	"strings"
	"testing"

	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

func TestGetValidArgsString(t *testing.T) {
	mdResp := plugin.GetMetadataResponse{
		Name:        "Example Plugin",
		Description: "This is an description of example plugin. 🍺",
		URL:         "https://example.com/notation/plugin",
		Version:     "1.0.0",
		Capabilities: []plugin.Capability{
			plugin.CapabilityEnvelopeGenerator,
			plugin.CapabilityTrustedIdentityVerifier,
			plugin.CapabilityRevocationCheckVerifier},
	}
	s := getValidArgsString(&mdResp)
	expected := "<generate-envelope|get-plugin-metadata|verify-signature|version>"
	if !strings.EqualFold(s, expected) {
		t.Errorf("Expected %s but found %s", expected, s)
	}
}

func TestGetValidArgs(t *testing.T) {
	tests := map[string]struct {
		caps []plugin.Capability
		args []string
	}{
		"sigGeneratorOnly": {
			caps: []plugin.Capability{plugin.CapabilitySignatureGenerator},
			args: []string{"describe-key", "generate-signature", "get-plugin-metadata", "version"},
		},
		"envGeneratorOnly": {
			caps: []plugin.Capability{plugin.CapabilityEnvelopeGenerator},
			args: []string{"generate-envelope", "get-plugin-metadata", "version"},
		},
		"verificationOnly": {
			caps: []plugin.Capability{plugin.CapabilityTrustedIdentityVerifier, plugin.CapabilityRevocationCheckVerifier},
			args: []string{"get-plugin-metadata", "verify-signature", "version"},
		},
		"envGenerator+verification": {
			caps: []plugin.Capability{plugin.CapabilityEnvelopeGenerator, plugin.CapabilityTrustedIdentityVerifier, plugin.CapabilityRevocationCheckVerifier},
			args: []string{"generate-envelope", "get-plugin-metadata", "verify-signature", "version"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mdResp := plugin.GetMetadataResponse{
				Name:         "Example Plugin",
				Description:  "This is an description of example plugin. 🍺",
				URL:          "https://example.com/notation/plugin",
				Version:      "1.0.0",
				Capabilities: test.caps,
			}
			s := getValidArgs(&mdResp)
			if !reflect.DeepEqual(s, test.args) {
				t.Errorf("Expected %s but found %s", test.args, s)
			}
		})
	}
}
