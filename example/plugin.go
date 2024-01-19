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

func (p *ExamplePlugin) DescribeKey(ctx context.Context, req *plugin.DescribeKeyRequest) (*plugin.DescribeKeyResponse, error) {
	return nil, plugin.NewUnsupportedError("DescribeKey operation is not implemented by example plugin")
}

func (p *ExamplePlugin) GenerateSignature(ctx context.Context, req *plugin.GenerateSignatureRequest) (*plugin.GenerateSignatureResponse, error) {
	return nil, plugin.NewUnsupportedError("GenerateSignature operation is not implemented by example plugin")
}

func (p *ExamplePlugin) GenerateEnvelope(ctx context.Context, req *plugin.GenerateEnvelopeRequest) (*plugin.GenerateEnvelopeResponse, error) {
	sig := "eyJwYXlsb2FkIjoiZXlKMFlYSm5aWFJCY25ScFptRmpkQ0k2ZXlKa2FXZGxjM1FpT2lKemFHRXlOVFk2Wm1VM1pUa3pNek16T1RVd05qQmpNbVkxWlRZelkyWXpObUV6" +
		"T0daaVlURXdNVGMyWmpFNE0ySTBNVFl6WVRVM09UUmxNRGd4WVRRNE1HRmlZbUUxWmlJc0ltMWxaR2xoVkhsd1pTSTZJbUZ3Y0d4cFkyRjBhVzl1TDNadVpDNWtiMk5yWlh" +
		"JdVpHbHpkSEpwWW5WMGFXOXVMbTFoYm1sbVpYTjBMbll5SzJwemIyNGlMQ0p6YVhwbElqbzVOREo5ZlEiLCJwcm90ZWN0ZWQiOiJleUpoYkdjaU9pSlFVekkxTmlJc0ltTn" +
		"lhWFFpT2xzaWFXOHVZMjVqWmk1dWIzUmhjbmt1YzJsbmJtbHVaMU5qYUdWdFpTSXNJbWx2TG1OdVkyWXVibTkwWVhKNUxuWmxjbWxtYVdOaGRHbHZibEJzZFdkcGJrMXBib" +
		"FpsY25OcGIyNGlMQ0pwYnk1amJtTm1MbTV2ZEdGeWVTNTJaWEpwWm1sallYUnBiMjVRYkhWbmFXNGlYU3dpWTNSNUlqb2lZWEJ3YkdsallYUnBiMjR2ZG01a0xtTnVZMll1" +
		"Ym05MFlYSjVMbkJoZVd4dllXUXVkakVyYW5OdmJpSXNJbWx2TG1OdVkyWXVibTkwWVhKNUxuTnBaMjVwYm1kVFkyaGxiV1VpT2lKdWIzUmhjbmt1ZURVd09TSXNJbWx2TG1" +
		"OdVkyWXVibTkwWVhKNUxuTnBaMjVwYm1kVWFXMWxJam9pTWpBeU15MHdNUzB4T1ZReE16b3dNem95TXkwd09Eb3dNQ0lzSW1sdkxtTnVZMll1Ym05MFlYSjVMblpsY21sbW" +
		"FXTmhkR2x2YmxCc2RXZHBiaUk2SW1sdkxtTnVZMll1Ym05MFlYSjVMbkJzZFdkcGJpNTFibWwwZEdWemRDNXRiMk5ySWl3aWFXOHVZMjVqWmk1dWIzUmhjbmt1ZG1WeWFXW" +
		"nBZMkYwYVc5dVVHeDFaMmx1VFdsdVZtVnljMmx2YmlJNklqRXVNQzR3TFdGc2NHaGhMbUpsZEdFaWZRIiwiaGVhZGVyIjp7Ing1YyI6WyJNSUlEVmpDQ0FqNmdBd0lCQWdJ" +
		"QlVUQU5CZ2txaGtpRzl3MEJBUXNGQURCYU1Rc3dDUVlEVlFRR0V3SlZVekVMTUFrR0ExVUVDQk1DVjBFeEVEQU9CZ05WQkFjVEIxTmxZWFIwYkdVeER6QU5CZ05WQkFvVEJ" +
		"rNXZkR0Z5ZVRFYk1Ca0dBMVVFQXhNU2QyRmlZbWwwTFc1bGRIZHZjbXR6TG1sdk1CNFhEVEl6TURFeE9UQTRNVGt3TjFvWERUTXpNREV4T1RBNE1Ua3dOMW93V2pFTE1Ba0" +
		"dBMVVFQmhNQ1ZWTXhDekFKQmdOVkJBZ1RBbGRCTVJBd0RnWURWUVFIRXdkVFpXRjBkR3hsTVE4d0RRWURWUVFLRXdaT2IzUmhjbmt4R3pBWkJnTlZCQU1URW5kaFltSnBkQ" +
		"zF1WlhSM2IzSnJjeTVwYnpDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTkhobFArU2lZN2hzR2xmMm1BRE96SlcvSjlzaXFNa2lRdlNPeDBP" +
		"U00yeXhldGZWUUwvYWJpNGlxQ1hNNndrU3h2aUJlTndJb1lFczR0aE1BOE5HRWJuS29Ya3R5aDl2bWlMQjFGVzdISHI0UUx3amdMemdXSktJUVR5MUptREJlY1haaDU2ZDB" +
		"mM3czWWoxSURUdmtJU2NYQ05JKzV2LzA4R1VRS2h5Qnd2N0ZxOU1ZcG8ybGZYU0k3VjMzQktLZGRYSXhQR1ZXd0tHdlBFMHNnMlZWN1dNODRaWkxkREt6Mm1xMFB0UFRIcl" +
		"N3ZzNobEsvbWpuK2JsZzNnc1lRNGg5LzdaNm5OYUY5WDBTZHlFU2w4NDFaV3J0TWhBT0ZwSXpMYno5ZXRlOE5SZDNiWUNSQklyNWdzY0hXVGY2bHlVZ3k0eHpzU3dNSFBzR" +
		"0xNNEErWjAwQ0F3RUFBYU1uTUNVd0RnWURWUjBQQVFIL0JBUURBZ2VBTUJNR0ExVWRKUVFNTUFvR0NDc0dBUVVGQndNRE1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQWJO" +
		"MEVydTU2dVRRU0MyOFpUZjhEN1Z5Q2tZcnJXTFlpSk1ZZE9LQnp6S1Y5bUthTTBPR0YydXlXd0RhUHhwOUtUZExYbUJwOUVGcTVTWFhBckZBK25SUzdLaW5EQWUyTzdBLzl" +
		"TdGQyWGpLaTkyN3JrQTJjajIzOWQ1bFJzaldYcUpYZjl2QU1WOWEyRmpVTS9pbjJFZXZscTdidmpGRTNsMjZWWENLdE9zOUVybWZ4ckwrNkVUUktTVllPT0cvclNIRnYvU0" +
		"IyTWxxRGc1UXNYQzlsWmp6TDUvWC9pb2UycVpLaHA2WDVEUHBhZDFxMVE0SXRLZFROKzJFWHlNeW9IbjFCSktOYmE3Q1VVdlhmMDNFSmViVC9JbStxb3pmRWtzSmVaSlVTb" +
		"FN1akFOVVBvQ3BzRVlHV1dReDVHK1ZpRzA1U3FzKzZwcEtydXQrUCtEVlBvIl0sImlvLmNuY2Yubm90YXJ5LnNpZ25pbmdBZ2VudCI6Ik5vdGF0aW9uLzEuMC4wIn0sInNp" +
		"Z25hdHVyZSI6ImlKdGh0cWJ6ME81bkZ1bzVaOW5SZGRFanlacDNSRy1LT1k2U1NCM3NjOEFnREJkVDVGanA5eWx0SW9xVGwtQkxaaHJHT0FGZU8wVF8xSlZzUGJaWk14ekp" +
		"xNGZiM2dQYUlQSXRyZW5ka3BpdDFtMlJhQjhmSzFEX0k2VnF1MV9yR2lZYXhEY05wYXFuMVRfaXN4cjRNVlJla2NMU05RbkczaU1kSjBrLUF0dGY4SmRDWEUwRVdLeUxCU3" +
		"RNVkFmbzBKMzlTaEZjd3lJTXZPMHZtMl9UUkRWYmNLb3ZwWTB2RnJmeUUycEZJQ2huSkVDbWl2SW1kS21CTUlXNzh2RXRONnFCcktza0kzSHpBOU4xWGp4R1k0R09BdTMwa" +
		"XF0TlJhbk82NW5aR25nMGxxcEpkMTViQXdVYXFqLUtEX0JBWklVVDlUMnFDZjJDT0Y5R0t2YzNOUSJ9Cg=="
	return &plugin.GenerateEnvelopeResponse{
		SignatureEnvelope:     []byte(sig),
		SignatureEnvelopeType: "application/jose+json",
		Annotations:           map[string]string{"manifestAnntnKey1": "value1"},
	}, nil
}

func (p *ExamplePlugin) VerifySignature(ctx context.Context, req *plugin.VerifySignatureRequest) (*plugin.VerifySignatureResponse, error) {
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

func (p *ExamplePlugin) GetMetadata(ctx context.Context, req *plugin.GetMetadataRequest) (*plugin.GetMetadataResponse, error) {
	return &plugin.GetMetadataResponse{
		Name:        "com.example.plugin",
		Description: "This is an description of example plugin",
		URL:         "https://example.com/notation/plugin",
		Version:     "1.0.0",
		Capabilities: []plugin.Capability{
			plugin.CapabilityEnvelopeGenerator,
			plugin.CapabilityTrustedIdentityVerifier,
			plugin.CapabilityRevocationCheckVerifier},
	}, nil
}
