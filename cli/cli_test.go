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
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/notaryproject/notation-plugin-framework-go/internal/mock"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

var cli, _ = New(mock.NewPlugin(false))
var errorCli, _ = New(mock.NewPlugin(true))

func TestNewWithLogger(t *testing.T) {
	_, err := NewWithLogger(nil, &discardLogger{})
	if err == nil {
		t.Errorf("NewWithLogger() expected error but not found")
	}
}

func TestMarshalResponse(t *testing.T) {
	res := plugin.GenerateEnvelopeResponse{
		SignatureEnvelope:     []byte("envelope"),
		Annotations:           map[string]string{"key": "value"},
		SignatureEnvelopeType: "envelopeType",
	}
	op, err := cli.marshalResponse(res, nil)
	if err != nil {
		t.Errorf("Error found in marshalResponse: %v", err)
	}

	expected := "{\"signatureEnvelope\":\"ZW52ZWxvcGU=\",\"signatureEnvelopeType\":\"envelopeType\",\"annotations\":{\"key\":\"value\"}}"
	if !strings.EqualFold("{\"signatureEnvelope\":\"ZW52ZWxvcGU=\",\"signatureEnvelopeType\":\"envelopeType\",\"annotations\":{\"key\":\"value\"}}", op) {
		t.Errorf("Not equal: \n expected: %s\n actual : %s", expected, op)
	}
}

func TestMarshalResponseError(t *testing.T) {
	_, err := cli.marshalResponse(nil, fmt.Errorf("expected error thrown"))
	assertErr(t, err, plugin.ErrorCodeGeneric)

	_, err = cli.marshalResponse(nil, plugin.NewValidationError("expected validation error thrown"))
	assertErr(t, err, plugin.ErrorCodeValidation)

	_, err = cli.marshalResponse(make(chan int), nil)
	assertErr(t, err, plugin.ErrorCodeGeneric)
}

func TestUnmarshalRequest(t *testing.T) {
	content := "{\"contractVersion\":\"1.0\",\"keyId\":\"someKeyId\",\"pluginConfig\":{\"pc1\":\"pk1\"}}"
	closer := setupReader(content)
	defer closer()

	var request plugin.DescribeKeyRequest
	if err := cli.unmarshalRequest(&request); err != nil {
		t.Errorf("unmarshalRequest() failed with error: %v", err)
	}

	if request.ContractVersion != "1.0" || request.KeyID != "someKeyId" || request.PluginConfig["pc1"] != "pk1" {
		t.Errorf("unmarshalRequest() returned incorrect struct")
	}
}

func TestUnmarshalRequestError(t *testing.T) {
	closer := setupReader("InvalidJson")
	defer closer()

	var request plugin.DescribeKeyRequest
	err := cli.unmarshalRequest(&request)
	if err == nil {
		t.Errorf("unmarshalRequest() expected error but not found")
	}

	plgErr, ok := err.(*plugin.Error)
	if !ok {
		t.Errorf("unmarshalRequest() expected error of type plugin.Error but found %s", reflect.TypeOf(err))
	}

	expectedErrStr := "{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"Input is not a valid JSON\"}"
	if plgErr.Error() != expectedErrStr {
		t.Errorf("unmarshalRequest() expected error string to be %s but found %s", expectedErrStr, plgErr.Error())
	}
}

func TestGetMetadataError(t *testing.T) {
	if os.Getenv("TEST_OS_EXIT") == "1" {
		ctx := context.Background()
		errorCli.Execute(ctx, []string{string(plugin.CommandGetMetadata)})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGetMetadataError")
	cmd.Env = append(os.Environ(), "TEST_OS_EXIT=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Errorf("process ran with err %v, want exit status 1", err)
}

func TestExecuteSuccess(t *testing.T) {
	sigGenCli, _ := New(mock.NewSigGeneratorPlugin(false))
	tests := map[string]struct {
		c  *CLI
		op string
	}{
		string(plugin.CommandGetMetadata): {
			c:  cli,
			op: "{\"name\":\"Example Plugin\",\"description\":\"This is an description of example plugin. 🍺\",\"version\":\"1.0.0\",\"url\":\"https://example.com/notation/plugin\",\"capabilities\":[\"SIGNATURE_VERIFIER.TRUSTED_IDENTITY\",\"SIGNATURE_VERIFIER.REVOCATION_CHECK\",\"SIGNATURE_GENERATOR.ENVELOPE\"]}",
		},
		string(plugin.Version): {
			c:  cli,
			op: "Example Plugin - This is an description of example plugin. 🍺\nVersion: 1.0.0 \n",
		},
		string(plugin.CommandGenerateEnvelope): {
			c:  cli,
			op: "{\"signatureEnvelope\":\"\",\"signatureEnvelopeType\":\"\",\"annotations\":{\"manifestAnntnKey1\":\"value1\"}}",
		},
		string(plugin.CommandVerifySignature): {
			c:  cli,
			op: "{\"verificationResults\":{\"SIGNATURE_VERIFIER.REVOCATION_CHECK\":{\"success\":true,\"reason\":\"Not revoked\"},\"SIGNATURE_VERIFIER.TRUSTED_IDENTITY\":{\"success\":true,\"reason\":\"Valid trusted Identity\"}},\"processedAttributes\":[]}",
		},
		string(plugin.CommandGenerateSignature): {
			c:  sigGenCli,
			op: "{\"keyId\":\"someKeyId\",\"signature\":\"YWJjZA==\",\"signingAlgorithm\":\"RSASSA-PSS-SHA-256\",\"certificateChain\":[\"YWJjZA==\",\"d3h5eg==\"]}",
		},
		string(plugin.CommandDescribeKey): {
			c:  sigGenCli,
			op: "{\"keyId\":\"someKeyId\",\"keySpec\":\"RSA-2048\"}",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			closer := setupReader("{}")
			defer closer()
			op := captureStdOut(func() {
				test.c.Execute(context.Background(), []string{"notation", name})
			})
			if op != test.op {
				t.Errorf("Execute() with '%s' args, expected '%s' but got '%s'", name, test.op, op)
			}
		})
	}
}

func setupReader(content string) func() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	os.Stdin = tmpfile
	return func() {
		os.Remove(tmpfile.Name())
		os.Stdin = oldStdin // Restore original Stdin

		tmpfile.Close()
	}
}

func captureStdOut(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

func assertErr(t *testing.T, err error, code plugin.ErrorCode) {
	if plgErr, ok := err.(*plugin.Error); ok {
		if reflect.DeepEqual(code, plgErr.ErrCode) {
			return
		}
		t.Errorf("mismatch in error code: \n expected: %s\n actual : %s", code, plgErr.ErrCode)
	}
	t.Errorf("expected error of type PluginError but found %s", reflect.TypeOf(err))
}
