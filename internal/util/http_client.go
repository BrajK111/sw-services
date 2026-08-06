package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var defaultHTTPClient = &http.Client{Timeout: 30 * time.Second}

// PostJSON POSTs body as JSON to rawURL (with optional query params) and
// unmarshals the response into out. A non-2xx status returns an error
// carrying the response body for diagnostics.
func PostJSON(rawURL string, query url.Values, body interface{}, out interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	if len(query) > 0 {
		rawURL = rawURL + "?" + query.Encode()
	}

	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("call %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response from %s: %w", rawURL, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s returned %d: %s", rawURL, resp.StatusCode, string(respBody))
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("unmarshal response from %s: %w", rawURL, err)
	}
	return nil
}

func JoinURL(host, path string) string {
	if len(host) > 0 && host[len(host)-1] != '/' {
		host += "/"
	}
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	return host + path
}
