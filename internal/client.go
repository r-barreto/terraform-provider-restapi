package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Call(endpoint string, path string, httpMethod string, requestBody string, id string, headers map[string]interface{}) (*string, error) {
	if strings.Contains(path, "{id}") {
		path = strings.ReplaceAll(path, "{id}", id)
	}
	if strings.Contains(requestBody, "{id}") {
		requestBody = strings.ReplaceAll(requestBody, "{id}", id)
	}

	httpRequest, err := http.NewRequest(httpMethod, fmt.Sprintf("%s%s", endpoint, path), strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	if requestBody != "" {
		httpRequest.Header.Set("Content-Type", "application/json")
	}

	for key, value := range headers {
		httpRequest.Header.Set(key, value.(string))
	}

	response, err := doRequest(httpRequest)
	if err != nil {
		return nil, err
	}

	responseString := string(response)

	return &responseString, nil
}

func doRequest(httpRequest *http.Request) ([]byte, error) {
	var client http.Client
	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", response.StatusCode, body)
	}

	return body, nil
}
