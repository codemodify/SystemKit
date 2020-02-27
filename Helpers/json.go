package Helpers

import (
	"encoding/json"
	"log"
)

// AsJSONString - takes an object and outputs a JSON formatted string
func AsJSONString(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		log.Println(err.Error())
	}
	return string(bytes)
}

// AsJSONStringWithIndentation - takes an object and outputs a JSON formatted string
func AsJSONStringWithIndentation(i interface{}) string {
	return AsJSONStringWithCustomIndentation(i, "\t")
}

// AsJSONStringWithCustomIndentation - takes an object and outputs a JSON formatted string
func AsJSONStringWithCustomIndentation(i interface{}, indentation string) string {
	bytes, err := json.MarshalIndent(i, "", indentation)
	if err != nil {
		log.Println(err.Error())
	}
	return string(bytes)
}
