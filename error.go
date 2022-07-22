package jsonrpc

import "errors"

const (
	// Invalid JSON was received by the server.
	// An error occurred on the server while parsing the JSON text.
	ParseError = -32700
	// The JSON sent is not a valid Request object.
	InvalidRequest = -32600
	// The method does not exist / is not available.
	MethodNotFound = -32601
	// Invalid method parameter(s).
	InvalidParams = -32602
	// Internal JSON-RPC error.
	InternalError = -32603
)

var ErrMethodExists = errors.New("method already registerd with rpc handler")

// errInvalidParamsType this error is returned from the Params struct
// during json unmarshaling of the request if the params is not
// of types nil, []interface{} or map[string]interface{}
var errInvalidParamsType error = errors.New("invalid prams type")

// internal errors used for building responses
var (
	errParseFailed = Error{
		Code:    ParseError,
		Message: "Parse error",
	}
	errBadRequest = Error{
		Code:    InvalidRequest,
		Message: "Invalid Request",
	}
	errBadParams = Error{
		Code:    InvalidParams,
		Message: "Invalid params format",
	}
	errMethodNotFound = Error{
		Code:    MethodNotFound,
		Message: "Method not found",
	}

	// this is a bit of a special caso
	// if for some reason json fails to marshal then this is a fallback
	// that can be written directly to the response writer
	errInternalBytes = []byte(`{"jsonrpc": "2.0", "error": {"code": -32603, "message": "Internal error"}, "id": null}`)
)
