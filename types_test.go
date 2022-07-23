package jsonrpc_test

import (
	"encoding/json"
	"testing"

	"github.com/indeedhat/json-rpc"
)

var paramsValidCases = []struct {
	name            string
	inputData       string
	expectedOutcome bool
}{
	{"string", `"this is a string"`, false},
	{"float", "3.14", false},
	{"int", "314", false},
	{"bool", "true", false},
	{"int slice", "[3, 1, 4]", true},
	{"string slice", `["three", "one", "four"]`, true},
	{"mixed slice", `["three", 1, 4.0]`, true},
	{"object of ints", `{"three": 3, "one": 1, "four": 4}`, true},
	{"object of strings", `{"three": "three", "one": "one", "four": "four"}`, true},
	{"mixed object", `{"three": "three", "one": 1, "four": 4.0}`, true},
}

func TestParamsValid(t *testing.T) {
	for _, testCase := range paramsValidCases {
		t.Run(testCase.name, func(t *testing.T) {
			var params jsonrpc.Params
			err := json.Unmarshal([]byte(testCase.inputData), &params)
			if err != nil {
				t.Fatal("failed to unmarshal data into params")
			}

			valid := params.Valid()
			if !valid && testCase.expectedOutcome {
				t.Fatal("false negative")
			} else if valid && !testCase.expectedOutcome {
				t.Fatal("false negative")
			}
		})
	}
}

type paramsUnmashalStruct struct {
	Name string
}

var paramsUnmarshalCases = []struct {
	name            string
	inputData       string
	outputTarget    any
	expectedOutcome bool
}{
	{"[]int->map[string]interface{}", "[3, 1, 4]", make(map[string]interface{}), false},
	{"[]int->[]string", "[3, 1, 4]", []string{}, false},
	{"[]string->[]int", `["three", "one", "four"]`, []int{}, false},
	{"[]interface{}->[]int", `["three", 1, 4.0]`, []int{}, false},
	{"[]interface{}->[]interface{}", `["three", 1, 4.0]`, []interface{}{}, true},
	{"[]int->struct", "[3, 1, 4]", paramsUnmashalStruct{}, false},
	{"object->invalidStruct", `{"name": 14}`, paramsUnmashalStruct{}, false},
	{"object->validStruct", `{"name": "Jimmy"}`, paramsUnmashalStruct{}, true},
}

func TestParamsUnmarshal(t *testing.T) {
	for _, testCase := range paramsUnmarshalCases {
		t.Run(testCase.name, func(t *testing.T) {
			var err error

			// im sure there is a better way of doing this but its late, im tired and this works
			switch testCase.outputTarget.(type) {
			case []int:
				var target []int
				err = json.Unmarshal([]byte(testCase.inputData), &target)

			case []string:
				var target []string
				err = json.Unmarshal([]byte(testCase.inputData), &target)

			case []interface{}:
				var target []interface{}
				err = json.Unmarshal([]byte(testCase.inputData), &target)

			case map[string]interface{}:
				var target = make(map[string]interface{})
				err = json.Unmarshal([]byte(testCase.inputData), &target)

			case paramsUnmashalStruct:
				var target paramsUnmashalStruct
				err = json.Unmarshal([]byte(testCase.inputData), &target)
			}

			if err != nil && testCase.expectedOutcome {
				t.Fatal("failed to unmarshal to valid target")
			} else if err == nil && !testCase.expectedOutcome {
				t.Fatalf("did not fail to unmarshal to invalid target: %v", testCase.outputTarget)
			}
		})
	}
}
