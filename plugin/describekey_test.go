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

func TestDescribeKeyRequest_Validate(t *testing.T) {
	reqs := []DescribeKeyRequest{
		getDescribeKeyRequest(ContractVersion, "someKeyId"),
		{
			ContractVersion: ContractVersion,
			KeyID:           "someKeyId",
			PluginConfig:    map[string]string{"someKey": "someValue"}},
	}

	for _, req := range reqs {
		if err := req.Validate(); err != nil {
			t.Errorf("VerifySignatureRequest#Validate failed with error: %+v", err)
		}
	}
}

func TestDescribeKeyRequest_Validate_Error(t *testing.T) {
	testCases := []struct {
		name string
		req  DescribeKeyRequest
	}{
		{name: "contractVersion", req: getDescribeKeyRequest("", "someKeyId")},
		{name: "keyId", req: getDescribeKeyRequest(ContractVersion, "")},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			if err := testcase.req.Validate(); err != nil {
				expMsg := fmt.Sprintf("{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"%s cannot be empty\"}", testcase.name)
				if err.Error() != expMsg {
					t.Errorf("expected error message '%s' but got '%s'", expMsg, err.Error())
				}
			} else {
				t.Error("DescribeKeyRequest#Validate didn't returned error")
			}
		})
	}
}

func TestDescribeKeyRequest_Command(t *testing.T) {
	req := getDescribeKeyRequest(ContractVersion, "someKeyId")
	if cmd := req.Command(); cmd != CommandDescribeKey {
		t.Errorf("DescribeKeyRequest#Command, expected %s but returned %s", CommandDescribeKey, cmd)
	}
}

func getDescribeKeyRequest(cv, kid string) DescribeKeyRequest {
	return DescribeKeyRequest{
		ContractVersion: cv,
		KeyID:           kid,
	}
}
