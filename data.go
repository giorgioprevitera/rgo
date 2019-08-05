// Package main provides ...
package main

type thing struct {
	Kind string                 `json:"kind"`
	Data map[string]interface{} `json:"data"`
}

type listing struct {
	Children []thing `json:"children,omitempty"`
}
