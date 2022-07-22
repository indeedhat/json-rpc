package jsonrpc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Handler implements the http.Handler interface
type Handler struct {
	handlers map[string]HandlerFunc
}

// NewHandler creates a new http.Handler for dealing with JSON-RPC 2.0 reqests
func NewHandler() *Handler {
	return &Handler{
		handlers: make(map[string]HandlerFunc),
	}
}

// RegisterRpcMethod registers an rpc method with the handler
// it will error out if the method name is already in use
func (h *Handler) RegisterRpcMethod(method string, handler HandlerFunc) error {
	if _, ok := h.handlers[method]; ok {
		return ErrMethodExists
	}

	h.handlers[method] = handler
	return nil
}

var _ http.Handler = (*Handler)(nil)

// ServeHTTP conforms to the http.Handler interfoce
func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	bodyData, err := ioutil.ReadAll(req.Body)

	if err != nil {
		rw.Write(marshal(buildResponse(ParseError, nil)))
		return
	}

	if !isBatchRequest(bodyData) {
		response := h.handleRequest(bodyData)
		if response != nil {
			rw.Write(marshal(response))
		}

		return
	}

	// TODO: this is a shitty solution, but it is a solution
	// there is probably something that can be done with a custom decoder
	var (
		responses   []Response
		tmpRequests []interface{}
	)

	// validity already checked
	_ = json.Unmarshal(bodyData, &tmpRequests)
	for _, tmpReq := range tmpRequests {
		bytes, _ := json.Marshal(tmpReq)
		response := h.handleRequest(bytes)

		if response != nil {
			responses = append(responses)
		}
	}

	rw.Write(marshal(responses))
}

// handleRequest handles each individual rpc request within the batch
func (h *Handler) handleRequest(body []byte) *Response {
	if !json.Valid(body) {
		return buildResponse(errParseFailed, nil)
	}

	var (
		req Request
		err = json.Unmarshal(body, &req)
	)

	if err != nil {
		if errors.Is(err, errInvalidParamsType) {
			return buildResponse(errBadParams, req.ID)
		}

		return buildResponse(errBadRequest, req.ID)
	}

	methodHandler, ok := h.handlers[req.Method]
	if !ok {
		return buildResponse(errMethodNotFound, req.ID)
	}

	return buildResponse(methodHandler(&req), req.ID)
}

// isBatchRequest checks for the furs none whitespace character in the body
// to see if it is an '[' (indicating that the json string is for a batch request)
func isBatchRequest(body []byte) bool {
	for _, char := range body {
		if char == ' ' || char == '\t' || char == '\r' || char == '\n' {
			continue
		}

		// this is the first none whitespace character, its the only
		// one we care about
		return char == '['
	}

	return false
}

// buildResponse builds a response object out of unknown input
func buildResponse(response any, id *any) *Response {
	// if no id is present then this is a notification request
	// and the server MUST not respond
	if id == nil {
		return nil
	}

	if err, ok := response.(Error); ok {
		return &Response{
			JsonRpc: "2.0",
			Error:   &err,
			ID:      id,
		}
	}

	return &Response{
		JsonRpc: "2.0",
		Result:  &response,
		ID:      id,
	}
}

// marshal the response object into a byte array for writing
func marshal(response any) []byte {
	data, err := json.Marshal(response)
	if err == nil {
		return data
	}

	return errInternalBytes
}
