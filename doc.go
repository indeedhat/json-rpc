// JSON-RPC 2.0 Handler
// This library provides a light weight net/http complient handler for consuming JSON-RPC 2.0 requests
// It is designed to allow you to add a json rpc endpoint to an existing HTTP server
//
// Spec complience
// This library should be fully spec complient as of the 2.0 specification, if you find something that
// i have missed please let me know and i will fix it, or submit a PR.
//
// Example
//
// 	handler := NewHandler()
//
// 	handler.RegisterRpcMethod("get.string", func(r *Request) any {
// 		return "string"
// 	})
//
// 	handler.RegisterRpcMethod("get.array", func(r *Request) any {
// 		return []string{"array", "of", "strings"}
// 	})
//
//  http.Handle("rpc", handler)
package jsonrpc
