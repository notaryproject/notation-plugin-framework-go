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
	"testing"
)

func TestNewError(t *testing.T) {
	msg := "someMSg"
	err := NewError(ErrorCodeValidation, msg)
	if err != nil {
		if ErrorCodeValidation != err.ErrCode {
			t.Errorf("NewError, expected errorCode '%s' but found'%s'", ErrorCodeValidation, err.ErrCode)
		}

		if msg != err.Message {
			t.Errorf("NewError, expected message'%s' but found '%s'", msg, err.Message)
		}

		if err.Metadata != nil {
			t.Errorf("NewError, expected metadata to be nil but found '%s'", err.Metadata)
		}

		expError := "{\"errorCode\":\"VALIDATION_ERROR\",\"errorMessage\":\"someMSg\"}"
		if expError != err.Error() {
			t.Errorf("NewError#Error, expected error to be '%s' but found '%s'", expError, err.Error())

		}

	} else {
		t.Error("NewError didn't return an error")
	}
}

func TestErrorCodes(t *testing.T) {
	testCases := []struct {
		err     *Error
		errCode ErrorCode
	}{
		{err: NewValidationError(""), errCode: ErrorCodeValidation},
		{err: NewValidationErrorf("%s", ""), errCode: ErrorCodeValidation},
		{err: NewUnsupportedError(""), errCode: ErrorCodeValidation},
		{err: NewGenericError(""), errCode: ErrorCodeGeneric},
		{err: NewGenericErrorf("%s", ""), errCode: ErrorCodeGeneric},
		{err: NewJSONParsingError(""), errCode: ErrorCodeValidation},
		{err: NewUnsupportedContractVersionError(""), errCode: ErrorCodeUnsupportedContractVersion},
	}
	for _, test := range testCases {
		if test.errCode != test.err.ErrCode {
			t.Errorf("Expected errorCode %s but found %s", test.errCode, test.err.ErrCode)
		}
	}
}
