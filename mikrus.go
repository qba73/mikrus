// Package mikrus is a client library for Mikrus VPS provider.
package mikrus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Clients represents Mikrus client.
type Client struct {
	apiKey     string
	serverID   string
	URL        string
	HTTPClient *http.Client
}

// New creates and returns new Mikrus client.
func New(apiKey, srvID string) Client {
	return Client{
		apiKey:   apiKey,
		serverID: srvID,
		URL:      "https://api.mikr.us",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Server represents information about Mirkus server.
type Server struct {
	ServerID       string `json:"server_id"`
	ServerName     string `json:"server_name,omitempty"`
	Expires        string `json:"expires"`
	ExpiresCytrus  string `json:"expires_cytrus,omitempty"`
	ExpiresStorage string `json:"expires_storage,omitempty"`
	ParamRam       string `json:"param_ram"`
	ParamDisk      string `json:"param_disk"`
	LastLogPanel   string `json:"lastlog_panel"`
	MikrusPro      string `json:"mikrus_pro"`
}

// Info returns information about server associated with the API Key and ServerID.
func (c *Client) Info(ctx context.Context) (Server, error) {
	res := Server{}
	val := url.Values{
		"key": []string{c.apiKey},
		"srv": []string{c.serverID},
	}
	if err := c.callAPI(ctx, "info", &res, []byte(val.Encode())); err != nil {
		return Server{}, err
	}
	return res, nil
}

// Servers returns short information about all servers associated
// with the API Key and ServerID.
func (c *Client) Servers(ctx context.Context) ([]Server, error) {
	servers := []Server{}
	val := url.Values{
		"key": []string{c.apiKey},
		"srv": []string{c.serverID},
	}
	if err := c.callAPI(ctx, "serwery", &servers, []byte(val.Encode())); err != nil {
		return []Server{}, err
	}
	return servers, nil
}

// Log represents a server log information.
type Log struct {
	ID          string `json:"id"`
	ServerID    string `json:"server_id"`
	Task        string `json:"task"`
	WhenCreated string `json:"when_created"`
	WhenDone    string `json:"when_done"`
	Output      string `json:"output"`
}

// Logs returns lats 10 log entries from the server associated
// with the API Key and ServerID.
func (c *Client) Logs(ctx context.Context) ([]Log, error) {
	logs := []Log{}
	val := url.Values{
		"key": []string{c.apiKey},
		"srv": []string{c.serverID},
	}
	if err := c.callAPI(ctx, "logs", &logs, []byte(val.Encode())); err != nil {
		return []Log{}, err
	}
	return logs, nil
}

func (c *Client) callAPI(ctx context.Context, verb string, res any, data []byte) error {
	requestURL := c.URL + "/" + verb
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	resp.Body.Close()
	respString := string(respBytes)
	resp.Body = io.NopCloser(strings.NewReader(respString))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status %d: %q", resp.StatusCode, respString)
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("decoding error for %q: %w", respString, err)
	}
	return nil
}
