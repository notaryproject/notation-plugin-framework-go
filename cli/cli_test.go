package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

func TestMarshalResponse(t *testing.T) {
	ctx := context.Background()
	res := plugin.GenerateEnvelopeResponse{
		SignatureEnvelope:     []byte("envelope"),
		Annotations:           map[string]string{"key": "value"},
		SignatureEnvelopeType: "envelopeType",
	}
	op, err := marshalResponse(ctx, res, nil)
	if err != nil {
		t.Errorf("Error found in marshalResponse: %v", err)
	}

	expected := "{\"signatureEnvelope\":\"ZW52ZWxvcGU=\",\"signatureEnvelopeType\":\"envelopeType\",\"annotations\":{\"key\":\"value\"}}"
	if !strings.EqualFold("{\"signatureEnvelope\":\"ZW52ZWxvcGU=\",\"signatureEnvelopeType\":\"envelopeType\",\"annotations\":{\"key\":\"value\"}}", op) {
		t.Errorf("Not equal: \n expected: %s\n actual : %s", expected, op)
	}
}

func TestMarshalResponseError(t *testing.T) {
	ctx := context.Background()

	_, err := marshalResponse(ctx, nil, fmt.Errorf("expected error thrown"))
	assertErr(t, err, plugin.CodeGeneric)

	_, err = marshalResponse(ctx, nil, plugin.NewValidationError("expected validation error thrown"))
	assertErr(t, err, plugin.CodeValidation)

	_, err = marshalResponse(ctx, make(chan int), nil)
	assertErr(t, err, plugin.CodeGeneric)
}

func TestUnmarshalRequest(t *testing.T) {
	ctx := context.Background()
	content := "{\"contractVersion\":\"1.0\",\"keyId\":\"someKeyId\",\"pluginConfig\":{\"pc1\":\"pk1\"}}"
	closer := setupReader(content)
	defer closer()

	var request plugin.DescribeKeyRequest
	if err := unmarshalRequest(ctx, &request); err != nil {
		t.Errorf("unmarshalRequest() failed with error: %v", err)
	}

	if request.ContractVersion != "1.0" || request.KeyID != "someKeyId" || request.PluginConfig["pc1"] != "pk1" {
		t.Errorf("unmarshalRequest() returned incorrect struct")
	}
}

func TestUnmarshalRequestError(t *testing.T) {
	ctx := context.Background()
	closer := setupReader("InvalidJson")
	defer closer()

	var request plugin.DescribeKeyRequest
	err := unmarshalRequest(ctx, &request)
	if err == nil {
		t.Fatalf("unmarshalRequest() expected error but not found")
	}

	plgErr, ok := err.(*plugin.Error)
	if !ok {
		t.Fatalf("unmarshalRequest() expected error of type plugin.Error but found %s", reflect.TypeOf(err))
	}

	expectedErrStr := "{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"Input is not a valid JSON\"}"
	if plgErr.Error() != expectedErrStr {
		t.Fatalf("unmarshalRequest() expected error string to be %s but found %s", expectedErrStr, plgErr.Error())
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

func assertErr(t *testing.T, err error, code plugin.Code) {
	if plgErr, ok := err.(*plugin.Error); ok {
		if reflect.DeepEqual(code, plgErr.ErrCode) {
			return
		}
		t.Errorf("mismatch in error code: \n expected: %s\n actual : %s", code, plgErr.ErrCode)
	}

	t.Errorf("expected error of type PluginError but found %s", reflect.TypeOf(err))
}
