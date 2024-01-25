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

package e2e

import (
	"encoding/json"
	"flag"
	"io"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var sigGenPluginPath = flag.String("sig_gen_plugin", "./bin/signaturegenerator/com.example.plugin", "dir of package containing embedded files")
var envGenPluginPath = flag.String("env_gen_plugin", "./bin/envelopegenerator/com.example.plugin", "dir of package containing embedded files")

func TestSuccess(t *testing.T) {
	tests := map[string]struct {
		pluginPath     *string
		stdin          string
		expectedStdout string
	}{
		"generate-envelope": {
			pluginPath:     envGenPluginPath,
			stdin:          "{\"contractVersion\":\"1.0\",\"keyId\":\"arn:aws:signer:us-west-2:951584113157:/signing-profiles/ECR\",\"payloadType\":\"application/vnd.cncf.notary.payload.v1+json\",\"signatureEnvelopeType\":\"application/jose+json\",\"payload\":\"eyJ0YXJnZXRBcnRpZmFjdCI6eyJtZWRpYVR5cGUiOiJhcHBsaWNhdGlvbi92bmQuZG9ja2VyLmRpc3RyaWJ1dGlvbi5tYW5pZmVzdC5saXN0LnYyK2pzb24iLCJkaWdlc3QiOiJzaGEyNTY6ZGEyN2I3NDAwOGJmYTQ5YTYyNWZhODVmNjZkMTJhMmY3YzE3ZGM0OTY4ZmJhZTZjOWRiNmU2N2ZkZDRmMjM4MiIsInNpemUiOjY4M319\"}",
			expectedStdout: "{\"signatureEnvelope\":\"ZXlKd1lYbHNiMkZrSWpvaVpYbEtNRmxZU201YVdGSkNZMjVTY0ZwdFJtcGtRMGsyWlhsS2EyRlhaR3hqTTFGcFQybEtlbUZIUlhsT1ZGazJXbTFWTTFwVWEzcE5lazE2VDFSVmQwNXFRbXBOYlZreFdsUlplbGt5V1hwT2JVVjZUMGRhYVZsVVJYZE5WR015V21wRk5FMHlTVEJOVkZsNldWUlZNMDlVVW14TlJHZDRXVlJSTkUxSFJtbFpiVVV4V21sSmMwbHRNV3hhUjJ4b1ZraHNkMXBUU1RaSmJVWjNZMGQ0Y0ZreVJqQmhWemwxVEROYWRWcEROV3RpTWs1eVdsaEpkVnBIYkhwa1NFcHdXVzVXTUdGWE9YVk1iVEZvWW0xc2JWcFlUakJNYmxsNVN6SndlbUl5TkdsTVEwcDZZVmh3YkVscWJ6Vk9SRW81WmxFaUxDSndjbTkwWldOMFpXUWlPaUpsZVVwb1lrZGphVTlwU2xGVmVra3hUbWxKYzBsdFRubGhXRkZwVDJ4emFXRlhPSFZaTWpWcVdtazFkV0l6VW1oamJtdDFZekpzYm1KdGJIVmFNVTVxWVVkV2RGcFRTWE5KYld4MlRHMU9kVmt5V1hWaWJUa3dXVmhLTlV4dVdteGpiV3h0WVZkT2FHUkhiSFppYkVKelpGZGtjR0pyTVhCaWJGcHNZMjVPY0dJeU5HbE1RMHB3WW5rMWFtSnRUbTFNYlRWMlpFZEdlV1ZUTlRKYVdFcHdXbTFzYWxsWVVuQmlNalZSWWtoV2JtRlhOR2xZVTNkcFdUTlNOVWxxYjJsWldFSjNZa2RzYWxsWVVuQmlNalIyWkcwMWEweHRUblZaTWxsMVltMDVNRmxZU2pWTWJrSm9aVmQ0ZGxsWFVYVmtha1Z5WVc1T2RtSnBTWE5KYld4MlRHMU9kVmt5V1hWaWJUa3dXVmhLTlV4dVRuQmFNalZ3WW0xa1ZGa3lhR3hpVjFWcFQybEtkV0l6VW1oamJtdDFaVVJWZDA5VFNYTkpiV3gyVEcxT2RWa3lXWFZpYlRrd1dWaEtOVXh1VG5CYU1qVndZbTFrVldGWE1XeEphbTlwVFdwQmVVMTVNSGROVXpCNFQxWlJlRTE2YjNkTmVtOTVUWGt3ZDA5RWIzZE5RMGx6U1cxc2RreHRUblZaTWxsMVltMDVNRmxZU2pWTWJscHNZMjFzYldGWFRtaGtSMngyWW14Q2MyUlhaSEJpYVVrMlNXMXNka3h0VG5WWk1sbDFZbTA1TUZsWVNqVk1ia0p6WkZka2NHSnBOVEZpYld3d1pFZFdlbVJETlhSaU1rNXlTV2wzYVdGWE9IVlpNalZxV21rMWRXSXpVbWhqYm10MVpHMVdlV0ZYV25CWk1rWXdZVmM1ZFZWSGVERmFNbXgxVkZkc2RWWnRWbmxqTW14MlltbEpOa2xxUlhWTlF6UjNURmRHYzJOSGFHaE1iVXBzWkVkRmFXWlJJaXdpYUdWaFpHVnlJanA3SW5nMVl5STZXeUpOU1VsRVZtcERRMEZxTm1kQmQwbENRV2RKUWxWVVFVNUNaMnR4YUd0cFJ6bDNNRUpCVVhOR1FVUkNZVTFSYzNkRFVWbEVWbEZSUjBWM1NsWlZla1ZNVFVGclIwRXhWVVZEUWsxRFZqQkZlRVZFUVU5Q1owNVdRa0ZqVkVJeFRteFpXRkl3WWtkVmVFUjZRVTVDWjA1V1FrRnZWRUpyTlhaa1IwWjVaVlJGWWsxQ2EwZEJNVlZGUVhoTlUyUXlSbWxaYld3d1RGYzFiR1JJWkhaamJYUjZURzFzZGsxQ05GaEVWRWw2VFVSRmVFOVVRVFJOVkd0M1RqRnZXRVJVVFhwTlJFVjRUMVJCTkUxVWEzZE9NVzkzVjJwRlRFMUJhMGRCTVZWRlFtaE5RMVpXVFhoRGVrRktRbWRPVmtKQloxUkJiR1JDVFZKQmQwUm5XVVJXVVZGSVJYZGtWRnBYUmpCa1IzaHNUVkU0ZDBSUldVUldVVkZMUlhkYVQySXpVbWhqYm10NFIzcEJXa0puVGxaQ1FVMVVSVzVrYUZsdFNuQmtRekYxV2xoU00ySXpTbkpqZVRWd1lucERRMEZUU1hkRVVWbEtTMjlhU1doMlkwNUJVVVZDUWxGQlJHZG5SVkJCUkVORFFWRnZRMmRuUlVKQlRraG9iRkFyVTJsWk4yaHpSMnhtTW0xQlJFOTZTbGN2U2psemFYRk5hMmxSZGxOUGVEQlBVMDB5ZVhobGRHWldVVXd2WVdKcE5HbHhRMWhOTm5kclUzaDJhVUpsVG5kSmIxbEZjelIwYUUxQk9FNUhSV0p1UzI5WWEzUjVhRGwyYldsTVFqRkdWemRJU0hJMFVVeDNhbWRNZW1kWFNrdEpVVlI1TVVwdFJFSmxZMWhhYURVMlpEQm1NM2N6V1dveFNVUlVkbXRKVTJOWVEwNUpLelYyTHpBNFIxVlJTMmg1UW5kMk4wWnhPVTFaY0c4eWJHWllVMGszVmpNelFrdExaR1JZU1hoUVIxWlhkMHRIZGxCRk1ITm5NbFpXTjFkTk9EUmFXa3hrUkV0Nk1tMXhNRkIwVUZSSWNsTjNaek5vYkVzdmJXcHVLMkpzWnpObmMxbFJOR2c1THpkYU5tNU9ZVVk1V0RCVFpIbEZVMnc0TkRGYVYzSjBUV2hCVDBad1NYcE1Zbm81WlhSbE9FNVNaRE5pV1VOU1FrbHlOV2R6WTBoWFZHWTJiSGxWWjNrMGVIcHpVM2ROU0ZCelIweE5ORUVyV2pBd1EwRjNSVUZCWVUxdVRVTlZkMFJuV1VSV1VqQlFRVkZJTDBKQlVVUkJaMlZCVFVKTlIwRXhWV1JLVVZGTlRVRnZSME5EYzBkQlVWVkdRbmROUkUxQk1FZERVM0ZIVTBsaU0wUlJSVUpEZDFWQlFUUkpRa0ZSUVdKT01FVnlkVFUyZFZSUlUwTXlPRnBVWmpoRU4xWjVRMnRaY25KWFRGbHBTazFaWkU5TFFucDZTMVk1YlV0aFRUQlBSMFl5ZFhsWGQwUmhVSGh3T1V0VVpFeFliVUp3T1VWR2NUVlRXRmhCY2taQksyNVNVemRMYVc1RVFXVXlUemRCTHpsVGRHUXlXR3BMYVRreU4zSnJRVEpqYWpJek9XUTFiRkp6YWxkWWNVcFlaamwyUVUxV09XRXlSbXBWVFM5cGJqSkZaWFpzY1RkaWRtcEdSVE5zTWpaV1dFTkxkRTl6T1VWeWJXWjRja3dyTmtWVVVrdFRWbGxQVDBjdmNsTklSbll2VTBJeVRXeHhSR2MxVVhOWVF6bHNXbXA2VERVdldDOXBiMlV5Y1ZwTGFIQTJXRFZFVUhCaFpERnhNVkUwU1hSTFpGUk9LekpGV0hsTmVXOUliakZDU2t0T1ltRTNRMVZWZGxobU1ETkZTbVZpVkM5SmJTdHhiM3BtUld0elNtVmFTbFZUYkZOMWFrRk9WVkJ2UTNCelJWbEhWMWRSZURWSEsxWnBSekExVTNGekt6WndjRXR5ZFhRclVDdEVWbEJ2SWwwc0ltbHZMbU51WTJZdWJtOTBZWEo1TG5OcFoyNXBibWRCWjJWdWRDSTZJazV2ZEdGMGFXOXVMekV1TUM0d0luMHNJbk5wWjI1aGRIVnlaU0k2SW1sS2RHaDBjV0o2TUU4MWJrWjFielZhT1c1U1pHUkZhbmxhY0ROU1J5MUxUMWsyVTFOQ00zTmpPRUZuUkVKa1ZEVkdhbkE1ZVd4MFNXOXhWR3d0UWt4YWFISkhUMEZHWlU4d1ZGOHhTbFp6VUdKYVdrMTRla3B4TkdaaU0yZFFZVWxRU1hSeVpXNWthM0JwZERGdE1sSmhRamhtU3pGRVgwazJWbkYxTVY5eVIybFpZWGhFWTA1d1lYRnVNVlJmYVhONGNqUk5WbEpsYTJOTVUwNVJia2N6YVUxa1NqQnJMVUYwZEdZNFNtUkRXRVV3UlZkTGVVeENVM1JOVmtGbWJ6QktNemxUYUVaamQzbEpUWFpQTUhadE1sOVVVa1JXWW1OTGIzWndXVEIyUm5KbWVVVXljRVpKUTJodVNrVkRiV2wyU1cxa1MyMUNUVWxYTnpoMlJYUk9ObkZDY2t0emEwa3pTSHBCT1U0eFdHcDRSMWswUjA5QmRUTXdhWEYwVGxKaGJrODJOVzVhUjI1bk1HeHhjRXBrTVRWaVFYZFZZWEZxTFV0RVgwSkJXa2xWVkRsVU1uRkRaakpEVDBZNVIwdDJZek5PVVNKOUNnPT0=\",\"signatureEnvelopeType\":\"application/jose+json\",\"annotations\":{\"manifestAnntnKey1\":\"value1\"}}",
		},
		"get-plugin-metadata": {
			pluginPath:     envGenPluginPath,
			stdin:          "{}",
			expectedStdout: "{\"name\":\"com.example.plugin\",\"description\":\"This is an description of example plugin\",\"version\":\"1.0.0\",\"url\":\"https://example.com/notation/plugin\",\"supportedContractVersions\":[\"1.0\"],\"capabilities\":[\"SIGNATURE_GENERATOR.ENVELOPE\",\"SIGNATURE_VERIFIER.TRUSTED_IDENTITY\",\"SIGNATURE_VERIFIER.REVOCATION_CHECK\"]}"},
		"verify-signature": {
			pluginPath:     envGenPluginPath,
			stdin:          "{\"contractVersion\":\"1.0\",\"signature\":{\"criticalAttributes\":{\"contentType\":\"someCT\",\"signingScheme\":\"someSigningScheme\"},\"unprocessedAttributes\":null,\"certificateChain\":[\"emFw\",\"em9w\"]},\"trustPolicy\":{\"trustedIdentities\":null,\"signatureVerification\":[\"SIGNATURE_GENERATOR.RAW\"]}}",
			expectedStdout: "{\"verificationResults\":{\"SIGNATURE_VERIFIER.REVOCATION_CHECK\":{\"success\":true,\"reason\":\"Not revoked\"},\"SIGNATURE_VERIFIER.TRUSTED_IDENTITY\":{\"success\":true,\"reason\":\"Valid trusted Identity\"}},\"processedAttributes\":[]}"},
		"version": {
			pluginPath:     envGenPluginPath,
			stdin:          "",
			expectedStdout: "com.example.plugin - This is an description of example plugin\nVersion: 1.0.0\n",
		},
		"generate-signature": {
			pluginPath:     sigGenPluginPath,
			stdin:          "{\"contractVersion\":\"1.0\",\"keyId\":\"someKeyId\",\"keySpec\":\"EC-384\",\"hashAlgorithm\":\"SHA-384\",\"payload\":\"em9w\"}",
			expectedStdout: "{\"keyId\":\"someKeyId\",\"signature\":\"Z2VuZXJhdGVkTW9ja1NpZ25hdHVyZQ==\",\"signingAlgorithm\":\"RSASSA-PSS-SHA-384\",\"certificateChain\":[\"bW9ja0NlcnQx\",\"bW9ja0NlcnQy\"]}",
		},
		"describe-key": {
			pluginPath:     sigGenPluginPath,
			stdin:          "{\"contractVersion\":\"1.0\",\"keyId\":\"someKeyId\"}",
			expectedStdout: "{\"keyId\":\"someKeyId\",\"keySpec\":\"RSA-3072\"}",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			validateSuccess(*test.pluginPath, name, test.stdin, test.expectedStdout, t)
		})
	}
}

func TestFailure(t *testing.T) {
	cmds := []string{"generate-envelope", "verify-signature", "get-plugin-metadata"}
	stdInputs := []string{"", "\n", "invalidjson", "🍺 ¢ §"}
	expectedValidationErr := "{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"Input is not a valid JSON\"}"
	for _, cmd := range cmds {
		for _, input := range stdInputs {
			t.Run(cmd+"_"+input, func(t *testing.T) {
				validateFailure(*envGenPluginPath, cmd, input, expectedValidationErr, t)
			})
		}
	}

	// get-plugin-metadata input has all of its keys as optional so ommiting it
	cmds = cmds[:len(cmds)-1]
	input := "{\"sad\":\"bad\"}"
	expectedValidationErr = "{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"Input is not a valid JSON: contractVersion cannot be empty\"}"
	for _, cmd := range cmds {
		t.Run(cmd+"_{}"+input, func(t *testing.T) {
			validateFailure(*envGenPluginPath, cmd, input, expectedValidationErr, t)
		})
	}

	cmds = []string{"generate-signature", "describe-key"}
	stdInputs = []string{"", "\n", "invalidjson", "🍺 ¢ §"}
	expectedValidationErr = "{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"Input is not a valid JSON\"}"
	for _, cmd := range cmds {
		for _, input := range stdInputs {
			t.Run(cmd+"_"+input, func(t *testing.T) {
				validateFailure(*sigGenPluginPath, cmd, input, expectedValidationErr, t)
			})
		}
	}
}

func execute(exe, arg, stdInput string, t *testing.T) (string, string, error) {
	cmd := exec.Command(exe, arg)

	if stdInput != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Errorf("something went wrong when trying to invoke example plugin by cli: %+v", err)
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, stdInput)
		}()
	}

	var stdOut strings.Builder
	var stdErr strings.Builder
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()

	return stdOut.String(), stdErr.String(), err
}

func validateSuccess(exe, arg, stdInput, expectedStdOut string, t *testing.T) {
	stdOut, stdErr, err := execute(exe, arg, stdInput, t)

	if err != nil {
		t.Logf("standard output: %s", stdOut)
		t.Logf("standard Err: %s", stdErr)
		t.Fatalf("'%s' command failed with error: %+v", arg, err)
	}

	if stdOut == "" {
		t.Errorf("For '%s' command's standard out must not be empty", arg)
	}

	if stdErr != "" {
		t.Errorf("For '%s' command's standard error must be empty", arg)
	}

	res, err := jsonEquals(stdOut, expectedStdOut)
	if err == nil {
		if !res {
			t.Errorf("For '%s' command, expected standard out to be '%s' but found '%s'", arg, expectedStdOut, stdOut)
		}
	} else {
		if expectedStdOut != stdOut {
			t.Errorf("For '%s' command, expected standard standard to be '%s' but found '%s'", arg, expectedStdOut, stdOut)
		}
	}

}

func validateFailure(exe, arg, stdInput, expectedStdError string, t *testing.T) {
	stdOut, stdErr, err := execute(exe, arg, stdInput, t)

	if err == nil {
		t.Fatalf("expected '%s' command fail with error but it didnt", arg)
		t.Logf("standard output: %s", stdOut)
		t.Logf("standard Err: %s", stdErr)
	}

	if stdErr == "" {
		t.Errorf("For '%s' command's standard error must not be empty", arg)
	}

	if stdOut != "" {
		t.Errorf("For '%s' command's standard out must be empty", arg)
	}

	if stdErr != expectedStdError {
		t.Errorf("For '%s' command, expected standard error to be '%s' but found '%s'", arg, expectedStdError, stdErr)
	}
}

// JSONEqual compares the JSON from two Readers.
func jsonEquals(x, y string) (bool, error) {
	var x1, y1 interface{}

	if err := json.Unmarshal([]byte(x), &x1); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(y), &y1); err != nil {
		return false, err
	}
	return reflect.DeepEqual(x1, y1), nil
}
