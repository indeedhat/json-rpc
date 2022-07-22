package jsonrpc

import "testing"

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
