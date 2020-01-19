package Helpers

import "encoding/json"

// AsJSONString - takes an object and outputs a JSON formatted string
func AsJSONString(i interface{}) string {
	bytes, _ := json.Marshal(i)
	return string(bytes)
}

// AsIndentedSONString - takes an object and outputs a JSON formatted string
func AsIndentedSONString(i interface{}) string {
	bytes, _ := json.MarshalIndent(i, "", "\t\t")
	return string(bytes)
}
