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

func TestGenerateSignatureRequest_Validate(t *testing.T) {
	reqs := []GenerateSignatureRequest{
		getGenerateSignatureRequest(ContractVersion, "someKeyId", string(KeySpecEC384), string(HashAlgorithmSHA384), []byte("zop")),
		{
			ContractVersion: ContractVersion,
			KeyID:           "someKeyId",
			KeySpec:         KeySpecEC384,
			Hash:            HashAlgorithmSHA384,
			Payload:         []byte("somePayload"),
			PluginConfig:    map[string]string{"key1": "value1"},
		},
	}

	for _, req := range reqs {
		if err := req.Validate(); err != nil {
			t.Errorf("VerifySignatureRequest#Validate failed with error: %+v", err)
		}
	}
}

func TestGenerateSignatureRequest_Validate_Error(t *testing.T) {
	testCases := []struct {
		name string
		req  GenerateSignatureRequest
	}{
		{name: "contractVersion", req: getGenerateSignatureRequest("", "someKeyId", string(KeySpecEC384), string(HashAlgorithmSHA384), []byte("zop"))},
		{name: "keyId", req: getGenerateSignatureRequest(ContractVersion, "", string(KeySpecEC384), string(HashAlgorithmSHA384), []byte("zop"))},
		{name: "keySpec", req: getGenerateSignatureRequest(ContractVersion, "someKeyId", "", string(HashAlgorithmSHA384), []byte("zop"))},
		{name: "hashAlgorithm", req: getGenerateSignatureRequest(ContractVersion, "someKid", string(KeySpecEC384), "", []byte("zop"))},
		{name: "payload", req: getGenerateSignatureRequest(ContractVersion, "someKeyId", string(KeySpecEC384), string(HashAlgorithmSHA384), []byte(""))},
		{name: "payload", req: getGenerateSignatureRequest(ContractVersion, "someKeyId", string(KeySpecEC384), string(HashAlgorithmSHA384), nil)},
		{name: "keySpec", req: getGenerateSignatureRequest(ContractVersion, "someKeyId", "", string(HashAlgorithmSHA384), []byte("zop"))},
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

func TestGenerateEnvelopeRequest_Validate(t *testing.T) {
	reqs := []GenerateEnvelopeRequest{
		getGenerateEnvelopeRequest(ContractVersion, "someKeyId", "someSET", "somePT", []byte("zop")),
		{
			ContractVersion:         ContractVersion,
			KeyID:                   "someKeyId",
			SignatureEnvelopeType:   "someType",
			ExpiryDurationInSeconds: 1,
			Payload:                 []byte("zap"),
			PayloadType:             ContractVersion,
			PluginConfig:            map[string]string{"key1": "value1"},
		},
	}

	for _, req := range reqs {
		if err := req.Validate(); err != nil {
			t.Errorf("GenerateEnvelopeRequest#Validate failed with error: %+v", err)
		}
	}
}

func TestGenerateEnvelopeRequest_Validate_Error(t *testing.T) {
	testCases := []struct {
		name string
		req  GenerateEnvelopeRequest
	}{
		{name: "contractVersion", req: getGenerateEnvelopeRequest("", "someKeyId", "someSET", "somePT", []byte("zop"))},
		{name: "keyId", req: getGenerateEnvelopeRequest(ContractVersion, "", "someSET", "somePT", []byte("zop"))},
		{name: "signatureEnvelopeType", req: getGenerateEnvelopeRequest(ContractVersion, "someKeyId", "", "somePT", []byte("zop"))},
		{name: "payloadType", req: getGenerateEnvelopeRequest(ContractVersion, "someKeyId", "someSET", "", []byte("zop"))},
		{name: "payload", req: getGenerateEnvelopeRequest(ContractVersion, "someKeyId", "someSET", "somePT", []byte(""))},
		{name: "payload", req: getGenerateEnvelopeRequest(ContractVersion, "someKeyId", "someSET", "somePT", nil)},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			if err := testcase.req.Validate(); err != nil {
				expMsg := fmt.Sprintf("{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"%s cannot be empty\"}", testcase.name)
				if err.Error() != expMsg {
					t.Errorf("expected error message '%s' but got '%s'", expMsg, err.Error())
				}
			} else {
				fmt.Println(testcase.req)
				t.Error("GenerateEnvelopeRequest#Validate didn't returned error")
			}
		})
	}
}

func getGenerateSignatureRequest(cv, kid, ks, ha string, pl []byte) GenerateSignatureRequest {
	return GenerateSignatureRequest{
		ContractVersion: cv,
		KeyID:           kid,
		KeySpec:         KeySpec(ks),
		Hash:            HashAlgorithm(ha),
		Payload:         pl,
	}
}

func getGenerateEnvelopeRequest(cv, kid, set, pt string, pl []byte) GenerateEnvelopeRequest {
	return GenerateEnvelopeRequest{
		ContractVersion:       cv,
		KeyID:                 kid,
		SignatureEnvelopeType: set,
		PayloadType:           pt,
		Payload:               pl,
	}
}
