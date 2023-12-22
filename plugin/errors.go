package plugin

import (
	"encoding/json"
	"fmt"
)

type Code string

const (
	CodeValidation                 Code = "VALIDATION_ERROR"
	CodeUnsupportedContractVersion Code = "UNSUPPORTED_CONTRACT_VERSION"
	CodeAccessDenied               Code = "ACCESS_DENIED"
	CodeThrottled                  Code = "THROTTLED"
	CodeGeneric                    Code = "ERROR"
)

const (
	ErrorMsgMalformedInput     string = "Input is not a valid JSON"
	ErrorMsgMalformedOutputFmt string = "Failed to generate response. Error: %s"
)

// Error is used when the signature associated is no longer
// valid.
type Error struct {
	ErrCode Code   `json:"errorCode"`
	Msg     string `json:"errorMessage"`
}

func NewError(code Code, msg string) *Error {
	return &Error{
		ErrCode: code,
		Msg:     msg,
	}
}

func NewGenericError(msg string) *Error {
	return NewError(CodeGeneric, msg)
}

func NewGenericErrorf(format string, msg string) *Error {
	return NewError(CodeGeneric, fmt.Sprintf(format, msg))
}

func NewUnsupportedError(msg string) *Error {
	return NewError(CodeValidation, msg+" is not supported")
}

func NewValidationError(msg string) *Error {
	return NewError(CodeValidation, msg)
}

func NewValidationErrorf(format string, msg string) *Error {
	return NewError(CodeValidation, fmt.Sprintf(format, msg))
}

func NewUnsupportedContractVersionError(version string) *Error {
	return NewError(CodeUnsupportedContractVersion, fmt.Sprintf("%q is not a supported notary plugin contract version", version))
}

func NewJSONParsingError(msg string) *Error {
	return NewValidationError(msg)
}

// Error returns the formatted error message.
func (e *Error) Error() string {
	op, err := json.Marshal(e)
	if err != nil {
		return "something went wrong"
	}
	return string(op)
}
