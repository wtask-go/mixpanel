package ingestion

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (*client) parseResponse(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("HTTP response is nil")
	}

	var (
		contentType = resp.Header.Get("Content-Type")
		err         error
	)

	switch {
	default:
		err = fmt.Errorf("unexpected response: %d %s, %s", resp.StatusCode, resp.Status, contentType)
	case resp.StatusCode == http.StatusOK && strings.Contains(contentType, "text/plain"):
		err = parsePlainText200(resp.Body)
	case resp.StatusCode == http.StatusOK && strings.Contains(contentType, "application/json"):
		err = parseJSON200(resp.Body)
	case (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) &&
		strings.Contains(contentType, "application/json"):
		err = parseJSONError(resp.Body)
	}

	return err
}

func parsePlainText200(body io.ReadCloser) error {
	defer func() {
		_ = body.Close()
	}()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read text/plain OK response: %s", err)
	}

	// API declared integer scheme for response, but as we known OpenAPI use float64
	response, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return fmt.Errorf("parse text/plain OK response: %s", err)
	}

	if response == 0 {
		return fmt.Errorf("request failed")
	}

	return nil
}

func parseJSON200(body io.ReadCloser) error {
	defer func() {
		_ = body.Close()
	}()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read application/json OK response: %s", err)
	}

	response := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("unmarshal application/json OK response: %s", err)
	}

	if response.Status == 0 {
		if response.Error == "" {
			response.Error = "details not provided"
		}

		return fmt.Errorf("request failed: %s", response.Error)
	}

	return nil
}

// parseJSONError parses generic error response
func parseJSONError(body io.ReadCloser) error {
	defer func() {
		_ = body.Close()
	}()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read application/json error response: %s", err)
	}

	content := struct {
		Status string `json:"status"` // most likely status value is "error" always
		Error  string `json:"error"`
	}{}
	if err := json.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("unmarshal application/json error response: %s", err)
	}

	if content.Error == "" {
		content.Error = "error details not provided"
	}

	return fmt.Errorf("request failed: %s", content.Error)
}
