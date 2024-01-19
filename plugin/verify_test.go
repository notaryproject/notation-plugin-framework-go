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

package plugin

import (
	"fmt"
	"testing"
)

var mockCertChain = [][]byte{[]byte("zap"), []byte("zop")}

func TestVerifySignatureRequest_Validate(t *testing.T) {
	reqs := []VerifySignatureRequest{
		getVerifySignatureRequest(ContractVersion, "someCT", "someSigningScheme", mockCertChain, []Capability{CapabilitySignatureGenerator}),
		{
			ContractVersion: "2.0",
			Signature: Signature{
				CriticalAttributes: CriticalAttributes{
					ContentType:          "someCT",
					SigningScheme:        "someSigningScheme",
					Expiry:               nil,
					AuthenticSigningTime: nil,
					ExtendedAttributes:   nil,
				},
				UnprocessedAttributes: []string{"upa1", "upa2"},
				CertificateChain:      mockCertChain,
			},
			TrustPolicy: TrustPolicy{
				TrustedIdentities:     []string{"trustedIdentity1", "trustedIdentity2"},
				SignatureVerification: []Capability{CapabilitySignatureGenerator, CapabilityRevocationCheckVerifier},
			},
			PluginConfig: map[string]string{"someKey": "someValue"},
		},
	}

	for _, req := range reqs {
		if err := req.Validate(); err != nil {
			t.Errorf("VerifySignatureRequest#Validate failed with error: %+v", err)
		}
	}
}

func TestVerifySignatureRequest_Validate_Error(t *testing.T) {
	reqWithoutSignature := getVerifySignatureRequest("1.0", "someCT", "someSigningScheme", mockCertChain, []Capability{CapabilitySignatureGenerator})
	reqWithoutSignature.Signature = Signature{}

	reqWithoutCriticalAttr := getVerifySignatureRequest("1.0", "someCT", "someSigningScheme", mockCertChain, []Capability{CapabilitySignatureGenerator})
	reqWithoutCriticalAttr.Signature.CriticalAttributes = CriticalAttributes{}

	testCases := []struct {
		name string
		req  VerifySignatureRequest
	}{
		{name: "contractVersion", req: getVerifySignatureRequest("", "someCT", "someSigningScheme", mockCertChain, []Capability{CapabilitySignatureGenerator})},
		{name: "signature's criticalAttributes's contentType", req: getVerifySignatureRequest(ContractVersion, "", "someSigningScheme", mockCertChain, []Capability{CapabilitySignatureGenerator})},
		{name: "signature's criticalAttributes's signingScheme", req: getVerifySignatureRequest(ContractVersion, "someCT", "", mockCertChain, []Capability{CapabilitySignatureGenerator})},
		{name: "signature's criticalAttributes's certificateChain", req: getVerifySignatureRequest(ContractVersion, "someCT", "someSigningScheme", [][]byte{}, []Capability{CapabilitySignatureGenerator})},
		{name: "signature's criticalAttributes's certificateChain", req: getVerifySignatureRequest(ContractVersion, "someCT", "someSigningScheme", nil, []Capability{CapabilitySignatureGenerator})},
		{name: "signature's trustPolicy", req: getVerifySignatureRequest(ContractVersion, "someCT", "someSigningScheme", mockCertChain, nil)},
		{name: "signature's trustPolicy's signatureVerification", req: getVerifySignatureRequest(ContractVersion, "someCT", "someSigningScheme", mockCertChain, []Capability{})},
		{name: "signature", req: reqWithoutSignature},
		{name: "signature's criticalAttributes", req: reqWithoutCriticalAttr},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			if err := testcase.req.Validate(); err != nil {
				expMsg := fmt.Sprintf("{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"%s cannot be empty\"}", testcase.name)
				if err.Error() != expMsg {
					t.Errorf("expected error message '%s' but got '%s'", expMsg, err.Error())
				}
			} else {
				t.Error("VerifySignatureRequest#Validate didn't returned error")
			}
		})
	}
}

func TestVerifySignatureRequest_Command(t *testing.T) {
	req := getVerifySignatureRequest(ContractVersion, "someCT", "someSigningScheme", mockCertChain, []Capability{CapabilitySignatureGenerator})
	if cmd := req.Command(); cmd != CommandVerifySignature {
		t.Errorf("DescribeKeyRequest#Command, expected %s but returned %s", CommandVerifySignature, cmd)
	}
}

func getVerifySignatureRequest(cv, ct, ss string, cc [][]byte, sv []Capability) VerifySignatureRequest {
	return VerifySignatureRequest{
		ContractVersion: cv,
		Signature: Signature{
			CriticalAttributes: CriticalAttributes{
				ContentType:   ct,
				SigningScheme: ss,
			},
			CertificateChain: cc,
		},
		TrustPolicy: TrustPolicy{
			SignatureVerification: sv,
		},
	}
}
