package ingestion

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/wtask-go/mixpanel/errs"
)

func (*client) parseResponse(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("%w: http response is nil", errs.ErrInvalidArgument)
	}

	var (
		contentType = resp.Header.Get("Content-Type")
		err         error
	)

	switch {
	default:
		err = fmt.Errorf(
			"%w: %d %s, %s",
			errs.ErrResponseInvalidContent,
			resp.StatusCode,
			resp.Status,
			contentType,
		)
	case resp.StatusCode == http.StatusOK && strings.Contains(contentType, "text/plain"):
		err = parsePlainText200(resp.Body)
	case resp.StatusCode == http.StatusOK && strings.Contains(contentType, "application/json"):
		err = parseJSON200(resp.Body)
	case (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) &&
		strings.Contains(contentType, "application/json"):
		err = parseJSONError(resp.Body)
		if err == nil {
			err = errs.ErrUnknown
		}

		err = fmt.Errorf("%w (%d %s)", err, resp.StatusCode, resp.Status)
	}

	return err
}

func parsePlainText200(body io.ReadCloser) error {
	defer func() {
		_ = body.Close()
	}()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("%w: read text/plain: %s", errs.ErrInvalidArgument, err)
	}

	// API declared integer scheme for response, but as we known OpenAPI use float64
	status, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return fmt.Errorf("%w: convert text/plain: %s", errs.ErrResponseInvalidContent, err)
	}

	if status == 0 {
		return fmt.Errorf("%w: request failed", errs.ErrResponseInvalidData)
	}

	return nil
}

func parseJSON200(body io.ReadCloser) error {
	defer func() {
		_ = body.Close()
	}()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("%w: read application/json: %s", errs.ErrInvalidArgument, err)
	}

	content := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}
	if err := json.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("%w: unmarshal json: %s", errs.ErrResponseInvalidContent, err)
	}

	if content.Status == 0 {
		if content.Error == "" {
			content.Error = "details not provided"
		}

		return fmt.Errorf("%w: request failed: %s", errs.ErrResponseInvalidData, content.Error)
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
		return fmt.Errorf("%w: read application/json error: %s", errs.ErrInvalidArgument, err)
	}

	content := struct {
		Status string `json:"status"` // most likely status value is "error" always 
		Error  string `json:"error"`
	}{}
	if err := json.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("%w: unmarshal json error: %s", errs.ErrResponseInvalidContent, err)
	}

	if content.Error == "" {
		content.Error = "error details not provided"
	}

	return fmt.Errorf("%w: %s", errs.ErrRequestFailed, content.Error)
}
