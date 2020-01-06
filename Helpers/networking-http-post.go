package Helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// DoPOSTWithJSON -
func DoPOSTWithJSON(endpointURL string, payload []byte) ([]byte, error) {
	request, err := http.NewRequest(
		"POST",
		endpointURL,
		bytes.NewBuffer(payload),
	)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)

	return data, err
}
