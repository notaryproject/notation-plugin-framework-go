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

package main

import (
	"context"

	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

type ExamplePlugin struct {
}

func NewExamplePlugin() (*ExamplePlugin, error) {
	return &ExamplePlugin{}, nil
}

func (p *ExamplePlugin) DescribeKey(_ context.Context, req *plugin.DescribeKeyRequest) (*plugin.DescribeKeyResponse, error) {
	return &plugin.DescribeKeyResponse{
		KeyID:   req.KeyID,
		KeySpec: plugin.KeySpecRSA3072,
	}, nil
}

func (p *ExamplePlugin) GenerateSignature(_ context.Context, req *plugin.GenerateSignatureRequest) (*plugin.GenerateSignatureResponse, error) {
	return &plugin.GenerateSignatureResponse{
		KeyID:            req.KeyID,
		Signature:        []byte("generatedMockSignature"),
		SigningAlgorithm: plugin.SignatureAlgorithmRSASSA_PSS_SHA384,
		CertificateChain: [][]byte{[]byte("mockCert1"), []byte("mockCert2")},
	}, nil
}

func (p *ExamplePlugin) GenerateEnvelope(_ context.Context, _ *plugin.GenerateEnvelopeRequest) (*plugin.GenerateEnvelopeResponse, error) {
	return nil, plugin.NewUnsupportedError("GenerateSignature operation is not implemented by example plugin")
}

func (p *ExamplePlugin) VerifySignature(_ context.Context, req *plugin.VerifySignatureRequest) (*plugin.VerifySignatureResponse, error) {
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

func (p *ExamplePlugin) GetMetadata(_ context.Context, _ *plugin.GetMetadataRequest) (*plugin.GetMetadataResponse, error) {
	return &plugin.GetMetadataResponse{
		SupportedContractVersions: []string{plugin.ContractVersion},
		Name:                      "com.example.plugin",
		Description:               "This is an description of example plugin",
		URL:                       "https://example.com/notation/plugin",
		Version:                   "1.0.0",
		Capabilities: []plugin.Capability{
			plugin.CapabilitySignatureGenerator,
			plugin.CapabilityTrustedIdentityVerifier,
			plugin.CapabilityRevocationCheckVerifier},
	}, nil
}
