package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
)

// NotificationResponse this type should be returned from HandlerFunc of the request
// being handled is a notification
type NotificationResponse struct {
}

// HandlerFunc specifies the handler type for registering new rpc hanlers
// any value can be returned but a value MUST be returned
// If the returned type is of Error then a standard JSON-RPC error response
// will be sent back to the client
//
// in the case of a notification request any none error responses will be discarded
type HandlerFunc func(*Request) any

type Params []byte

var _ json.Unmarshaler = (*Params)(nil)

// UnmarshalJSON provides custom logic for encoding/json to unmarshal data into the Params type
func (p *Params) UnmarshalJSON(data []byte) error {
	*p = data
	return nil
}

// Valid checks the validity of the params input (if they are an array or object)
func (p *Params) Valid() bool {
	if len(*p) == 0 {
		return true
	}

	var tmpData interface{}
	if err := json.Unmarshal(*p, &tmpData); err != nil {
		return false
	}

	switch tmpData.(type) {
	case map[string]interface{}, []interface{}:
		return true
	}

	return false
}

// Unmarshal unmarshals the params byte array into the given target structure
// if unmarshaling fails an rpc error is returned that can either directly be
// returned from the handler to the client or ignored in favour of a programmer
// defined error
func (p Params) Unmarshal(target any) error {
	err := json.Unmarshal(p, target)
	if err != nil {
		return Error{
			Code:    InvalidParams,
			Message: "Invalid Request Params",
			Data:    err,
		}
	}

	return nil
}

// Request represents an rpc request sent to the server
type Request struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      *any   `json:"id"`
	Ctx     context.Context

	// parsed is used as part of batch requesting to determine if
	parsed bool
}

// Response represents an rpc response sent back from the server
type Response struct {
	JsonRpc string `json:"jsonrpc"`
	Result  *any   `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
	ID      *any   `json:"id"`
}

// Error represents an error handling an rpc reqest
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var _ error = (*Error)(nil)

// Error this makes the Error type conform to the error interface
func (e Error) Error() string {
	if e.Data == nil {
		return fmt.Sprintf("rpc-error(%d): %s", e.Code, e.Message)
	}

	return fmt.Sprintf("rpc-error(%d): %s\n%v", e.Code, e.Message, e.Data)
}
