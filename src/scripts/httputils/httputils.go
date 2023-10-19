package httputils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func SendHTTPRequest(method, url, body string, headers map[string]string) (string, int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %v", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("error reading response: %v", err)
	}

	return string(respBody), resp.StatusCode, nil
}