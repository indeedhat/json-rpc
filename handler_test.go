package jsonrpc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var isBatchRequestCases = []struct {
	name            string
	input           string
	expectedOutcome bool
}{
	{"invalid json", "notjson", false},
	{"empty string", "", false},
	{"batch", `[{"jsonrpc": "2.0", "method": "exapmle"}]`, true},
	{"batch with whitespace", `  
 	           [{"jsonrpc": "2.0", "method": "exapmle"}]`, true},
	{"single", `{"jsonrpc": "2.0", "method": "exapmle"}`, false},
	{"single with whitespace", `  
 	           {"jsonrpc": "2.0", "method": "exapmle"}`, false},
}

func TestIsBatchRequest(t *testing.T) {
	for _, testCase := range isBatchRequestCases {
		t.Run(testCase.name, func(t *testing.T) {
			isBatch := isBatchRequest([]byte(testCase.input))
			if isBatch && !testCase.expectedOutcome {
				t.Fatal("false positive")
			} else if !isBatch && testCase.expectedOutcome {
				t.Fatal("false negative")
			}
		})
	}
}

var registerRpcMethodCases = []struct {
	name            string
	method          string
	expectedOutcome bool
}{
	{"unique 1", "test.unique", true},
	{"unique 2", "test.alsounique", true},
	{"repeated", "test.unique", false},
}

func TestRegisterRpcMethod(t *testing.T) {
	handler := NewHandler()

	for _, testCase := range registerRpcMethodCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := handler.RegisterRpcMethod(testCase.method, func(r *Request) any {
				return nil
			})

			if err == nil && !testCase.expectedOutcome {
				t.Fatal("false positive")
			} else if err != nil && testCase.expectedOutcome {
				t.Fatal("false negative")
			}
		})
	}
}

var serveHttpCases = []struct {
	name           string
	input          string
	expectedOutput string
}{
	{
		"empty body",
		"",
		`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`,
	},
	{
		"bad json",
		"{not json",
		`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`,
	},
	{
		"bad version",
		`{"jsonrpc":"1.0", "method":"get.string"}`,
		`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`,
	},
	{
		"bad params",
		`{"jsonrpc":"2.0", "method":"get.string", "params": "string", "id":1}`,
		`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params format"},"id":1}`,
	},
	{
		"bad params (no id)",
		`{"jsonrpc":"2.0", "method":"get.string", "params": "string"}`,
		``,
	},
	{
		"not found",
		`{"jsonrpc":"2.0", "method":"get.missing", "id":1}`,
		`{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1}`,
	},
	{
		"not found (no id)",
		`{"jsonrpc":"2.0", "method":"get.missing"}`,
		``,
	},
	{
		"get string",
		`{"jsonrpc":"2.0", "method":"get.string", "params": ["string"], "id": 2}`,
		`{"jsonrpc":"2.0","result":"string","id":2}`,
	},
	{
		"get array",
		`{"jsonrpc":"2.0", "method":"get.array", "params": {"object": "params"}, "id": 3}`,
		`{"jsonrpc":"2.0","result":["array","of","strings"],"id":3}`,
	},
	{
		"get object",
		`{"jsonrpc":"2.0", "method":"get.map", "params": {"object": "params"}, "id": "str-id"}`,
		`{"jsonrpc":"2.0","result":{"key1":"val1","key2":2},"id":"str-id"}`,
	},
	{
		"get object (no id)",
		`{"jsonrpc":"2.0", "method":"get.array", "params": {"object": "params"}}`,
		``,
	},
	{
		"get params (array)",
		`{"jsonrpc":"2.0", "method":"get.params", "params": ["array","params"], "id": "str-id"}`,
		`{"jsonrpc":"2.0","result":["array","params"],"id":"str-id"}`,
	},
	{
		"get params (object)",
		`{"jsonrpc":"2.0", "method":"get.params", "params": {"object": "params"}, "id": "str-id"}`,
		`{"jsonrpc":"2.0","result":{"object":"params"},"id":"str-id"}`,
	},
	{
		"batch request",
		`[
            "",
            "{not json",
            {"jsonrpc":"1.0", "method":"get.string"},
            {"jsonrpc":"2.0", "method":"get.string", "params": "string", "id":1},
            {"jsonrpc":"2.0", "method":"get.string", "params": "string"},
            {"jsonrpc":"2.0", "method":"get.missing", "id":1},
            {"jsonrpc":"2.0", "method":"get.missing"},
            {"jsonrpc":"2.0", "method":"get.string", "params": ["string"], "id": 2},
            {"jsonrpc":"2.0", "method":"get.array", "params": {"object": "params"}, "id": 3},
            {"jsonrpc":"2.0", "method":"get.map", "params": {"object": "params"}, "id": "str-id"},
            {"jsonrpc":"2.0", "method":"get.map", "params": {"object": "params"}},
            {"jsonrpc":"2.0", "method":"get.params", "params": ["array","params"], "id": "str-id"},
            {"jsonrpc":"2.0", "method":"get.params", "params": {"object": "params"}, "id": "str-id"}
        ]`,
		`[` +
			`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null},` +
			`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null},` +
			`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null},` +
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params format"},"id":1},` +
			`{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1},` +
			`{"jsonrpc":"2.0","result":"string","id":2},` +
			`{"jsonrpc":"2.0","result":["array","of","strings"],"id":3},` +
			`{"jsonrpc":"2.0","result":{"key1":"val1","key2":2},"id":"str-id"},` +
			`{"jsonrpc":"2.0","result":["array","params"],"id":"str-id"},` +
			`{"jsonrpc":"2.0","result":{"object":"params"},"id":"str-id"}` +
			`]`,
	},
}

func TestServeHttp(t *testing.T) {
	handler := NewHandler()

    http.Handle("rpc", handler)
	handler.RegisterRpcMethod("get.string", func(r *Request) any {
		return "string"
	})
	handler.RegisterRpcMethod("get.array", func(r *Request) any {
		return []string{"array", "of", "strings"}
	})
	handler.RegisterRpcMethod("get.map", func(r *Request) any {
		return map[string]interface{}{
			"key1": "val1",
			"key2": 2,
		}
	})
	handler.RegisterRpcMethod("get.params", func(r *Request) any {
		var params interface{}
		r.Params.Unmarshal(&params)
		return params
	})

	for _, testCase := range serveHttpCases {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.input)
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "", reader)
			if err != nil {
				t.Fatal("failed to create request")
			}

			handler.ServeHTTP(recorder, req)

			response := string(recorder.Body.Bytes())
			if response != testCase.expectedOutput {
				fmt.Println(testCase.expectedOutput)
				t.Fatalf("bad response: %s", response)
			}
		})
	}
}
