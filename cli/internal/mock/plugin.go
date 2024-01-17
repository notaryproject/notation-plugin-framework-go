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

// Package mock contains various mock structures required for testing.
package mock

import (
	"context"
	"fmt"

	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

// Mock plugin used only for testing.
type mockPlugin struct {
	fail bool
}

func NewPlugin(failPlugin bool) *mockPlugin {
	return &mockPlugin{fail: failPlugin}
}

func (p *mockPlugin) DescribeKey(ctx context.Context, req *plugin.DescribeKeyRequest) (*plugin.DescribeKeyResponse, error) {
	return nil, plugin.NewUnsupportedError("DescribeKey operation is not implemented by example plugin")
}

func (p *mockPlugin) GenerateSignature(ctx context.Context, req *plugin.GenerateSignatureRequest) (*plugin.GenerateSignatureResponse, error) {
	return nil, plugin.NewUnsupportedError("GenerateSignature operation is not implemented by example plugin")
}

func (p *mockPlugin) GenerateEnvelope(ctx context.Context, req *plugin.GenerateEnvelopeRequest) (*plugin.GenerateEnvelopeResponse, error) {
	if p.fail {
		return nil, fmt.Errorf("expected error")
	}
	return &plugin.GenerateEnvelopeResponse{
		SignatureEnvelope:     []byte(""),
		SignatureEnvelopeType: "application/jose+json",
		Annotations:           map[string]string{"manifestAnntnKey1": "value1"},
	}, nil
}

func (p *mockPlugin) VerifySignature(ctx context.Context, req *plugin.VerifySignatureRequest) (*plugin.VerifySignatureResponse, error) {
	if p.fail {
		return nil, fmt.Errorf("expected error")
	}
	upAttrs := req.Signature.UnprocessedAttributes
	pAttrs := make([]interface{}, len(upAttrs))
	for i := range upAttrs {
		pAttrs[i] = upAttrs[i]
	}

	return &plugin.VerifySignatureResponse{
		ProcessedAttributes: pAttrs,
		VerificationResults: map[plugin.Capability]*plugin.VerificationResult{
			plugin.CapabilityTrustedIdentityVerifier: {
				Success: true,
				Reason:  "Valid trusted Identity",
			},
			plugin.CapabilityRevocationCheckVerifier: {
				Success: true,
				Reason:  "Not revoked",
			},
		},
	}, nil
}

func (p *mockPlugin) GetMetadata(ctx context.Context, req *plugin.GetMetadataRequest) (*plugin.GetMetadataResponse, error) {
	if p.fail {
		return nil, fmt.Errorf("expected error")
	}
	return &plugin.GetMetadataResponse{
		Name:        "Example Plugin",
		Description: "This is an description of example plugin. 🍺",
		URL:         "https://example.com/notation/plugin",
		Version:     "1.0.0",
		Capabilities: []plugin.Capability{
			plugin.CapabilityEnvelopeGenerator,
			plugin.CapabilityTrustedIdentityVerifier,
			plugin.CapabilityRevocationCheckVerifier},
	}, nil
}
