package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/notaryproject/notation-plugin-framework-go/cli/internal/mock"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

var cli = New("mockCli", &mockPlugin{})

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
	assertErr(t, err, plugin.CodeGeneric)

	_, err = cli.marshalResponse(nil, plugin.NewValidationError("expected validation error thrown"))
	assertErr(t, err, plugin.CodeValidation)

	_, err = cli.marshalResponse(make(chan int), nil)
	assertErr(t, err, plugin.CodeGeneric)
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

func TestGetMetadata(t *testing.T) {
	pl := mock.NewPlugin(false)
	ctx := context.Background()

	cli.getMetadata(ctx, pl)
}

func TestGetMetadataError(t *testing.T) {
	if os.Getenv("TEST_OS_EXIT") == "1" {
		pl := mock.NewPlugin(true)
		ctx := context.Background()

		cli.getMetadata(ctx, pl)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGetMetadataError")
	cmd.Env = append(os.Environ(), "TEST_OS_EXIT=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
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

type mockPlugin struct {
}

func (p *mockPlugin) DescribeKey(_ context.Context, _ *plugin.DescribeKeyRequest) (*plugin.DescribeKeyResponse, error) {
	return nil, plugin.NewUnsupportedError("DescribeKey operation is not implemented by example plugin")
}

func (p *mockPlugin) GenerateSignature(_ context.Context, _ *plugin.GenerateSignatureRequest) (*plugin.GenerateSignatureResponse, error) {
	return nil, plugin.NewUnsupportedError("GenerateSignature operation is not implemented by example plugin")
}

func (p *mockPlugin) GenerateEnvelope(_ context.Context, _ *plugin.GenerateEnvelopeRequest) (*plugin.GenerateEnvelopeResponse, error) {
	return nil, plugin.NewUnsupportedError("GenerateEnvelope operation is not implemented by example plugin")
}

func (p *mockPlugin) VerifySignature(_ context.Context, _ *plugin.VerifySignatureRequest) (*plugin.VerifySignatureResponse, error) {
	return nil, plugin.NewUnsupportedError("VerifySignature operation is not implemented by example plugin")

}

func (p *mockPlugin) GetMetadata(_ context.Context, _ *plugin.GetMetadataRequest) (*plugin.GetMetadataResponse, error) {
	return nil, plugin.NewUnsupportedError("GetMetadata operation is not implemented by mock plugin")
}
