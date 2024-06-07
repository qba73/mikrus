// Package mikrus is a client library for Mikrus VPS provider.
package mikrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents Mikrus client.
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

// Info returns information about server associated with the API Key and ServerID.
func (c *Client) Info() (Server, error) {
	res := Server{}
	if err := c.callAPI("info", &res); err != nil {
		return Server{}, err
	}
	return res, nil
}

// Servers returns short information about all servers associated
// with the API Key and ServerID.
func (c *Client) Servers() (Servers, error) {
	servers := Servers{}
	if err := c.callAPI("serwery", &servers); err != nil {
		return Servers{}, err
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
func (c *Client) Logs() ([]Log, error) {
	logs := []Log{}
	if err := c.callAPI("logs", &logs); err != nil {
		return []Log{}, err
	}
	return logs, nil
}

func (c *Client) callAPI(verb string, res any) error {
	requestURL := c.URL + "/" + verb
	val := url.Values{
		"key": []string{c.apiKey},
		"srv": []string{c.serverID},
	}
	resp, err := c.HTTPClient.PostForm(requestURL, val)
	if err != nil {
		return err
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

// ServerShort represents short server description.
type ServerShort struct {
	ServerID   string `json:"server_id"`
	ServerName string `json:"server_name"`
	Expires    string `json:"expires"`
	ParamRam   string `json:"param_ram"`
	ParamDisk  string `json:"param_disk"`
}

const serversTemplate = `{{ range . }}
Server ID: {{ .ServerID }}
Server name: {{ .ServerName }}
Expiration date: {{ .Expires }}
RAM size: {{ .ParamRam }}
ParamDisk: {{ .ParamDisk }}
{{ end }}`

// Servers is a list of servers in a short form.
type Servers []ServerShort

// String implements stringer interface.
func (s Servers) String() string {
	out, err := render(serversTemplate, s)
	if err != nil {
		return fmt.Sprintln(err.Error())
	}
	return out
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

const serverTemplate = `ServerID: {{ .ServerID }}
Server name: {{ .ServerName }}
Expiration date: {{ .Expires }}
Cytrus expiration date: {{ .ExpiresCytrus }}
Storage expiration date: {{ .ExpiresStorage }}
RAM size: {{ .ParamRam }}
Disk size: {{ .ParamDisk }}
Last log time: {{ .LastLogPanel }}
Is Pro service: {{ .MikrusPro }}`

// String implements stringer interface.
func (s Server) String() string {
	out, err := render(serverTemplate, s)
	if err != nil {
		return fmt.Sprintln(err.Error())
	}
	return out
}

// render takes a template and a data value, and returns
// the string result of executing the template.
func render(templateName string, value any) (string, error) {
	tmpl, err := template.New("").Parse(templateName)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, value)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
