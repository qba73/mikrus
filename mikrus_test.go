package mikrus_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/mikrus"
)

func TestMikrusReturnsInformationAboutServer(t *testing.T) {
	t.Parallel()

	ts := newTestServer("/info", info, t)
	defer ts.Close()

	c := mikrus.New("dummyKey", "dummySrv")
	c.HTTPClient = ts.Client()
	c.URL = ts.URL
	got, err := c.Info()
	if err != nil {
		t.Fatal(err)
	}
	want := mikrus.Server{
		ServerID:     "j230",
		Expires:      "2026-06-08 00:00:00",
		ParamRam:     "1024",
		ParamDisk:    "10",
		LastLogPanel: "2024-06-05 10:02:55",
		MikrusPro:    "nie",
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMikrusReturnsListOfServers(t *testing.T) {
	t.Parallel()

	ts := newTestServer("/serwery", servers, t)
	defer ts.Close()

	c := mikrus.New("dummyAPIKey", "dummyServerID")
	c.HTTPClient = ts.Client()
	c.URL = ts.URL

	got, err := c.Servers()
	if err != nil {
		t.Fatal(err)
	}
	want := mikrus.Servers{
		{
			ServerID:  "a133",
			Expires:   "2025-06-05 00:00:00",
			ParamRam:  "1024",
			ParamDisk: "10",
		},
		{
			ServerID:  "j139",
			Expires:   "2026-06-08 00:00:00",
			ParamRam:  "1024",
			ParamDisk: "10",
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMikrusReturnsServerLogs(t *testing.T) {
	t.Parallel()

	ts := newTestServer("/logs", logs, t)
	defer ts.Close()

	c := mikrus.New("dummyAPIKey", "dummyServerID")
	c.HTTPClient = ts.Client()
	c.URL = ts.URL

	got, err := c.Logs()
	if err != nil {
		t.Fatal(err)
	}
	want := []mikrus.Log{
		{
			ID:          "3752",
			ServerID:    "j230",
			Task:        "kluczssh",
			WhenCreated: "2024-06-05 10:05:34",
			WhenDone:    "2024-06-05 10:06:01",
			Output:      "Wrzuciłem klucz SSH\n",
		},
		{
			ID:          "3751",
			ServerID:    "j230",
			Task:        "restart",
			WhenCreated: "2024-06-05 09:57:54",
			WhenDone:    "2024-06-05 09:58:07",
			Output:      "OK\n",
		},
		{
			ID:          "3748",
			ServerID:    "j230",
			Task:        "upgrade",
			WhenCreated: "2024-06-05 08:59:28",
			WhenDone:    "2024-06-05 09:00:04",
			Output:      "=== Aktualne parametry: 768 RAM / 10 DYSK\n2 / 20\nDodaje: +256MB RAM oraz +0GB dysku\nPo zmianie: 1024 MB / 10 GB\n[succes] GOTOWE!\n",
		},
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func newTestServer(path string, data []byte, t *testing.T) *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("want POST request, got %q", r.Method)
		}
		verifyURL(path, r.URL.EscapedPath(), t)
		if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
			t.Fatal(err)
		}
	}))
}

// verifyURL checks if provided URLs are equal.
func verifyURL(wantURL, path string, t *testing.T) {
	t.Helper()

	wantU, err := url.Parse(wantURL)
	if err != nil {
		t.Fatalf("error parsing URL %q, %v", wantURL, err)
	}
	gotU, err := url.Parse(path)
	if err != nil {
		t.Fatalf("error parsing URL %q, %v", wantURL, err)
	}

	if !cmp.Equal(wantU.Path, gotU.Path) {
		t.Fatalf(cmp.Diff(wantU.Path, gotU.Path))
	}

	wantQuery, err := url.ParseQuery(wantU.RawQuery)
	if err != nil {
		t.Fatal(err)
	}
	gotQuery, err := url.ParseQuery(gotU.RawQuery)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wantQuery, gotQuery) {
		t.Fatalf("URLs are not equal, \n%s", cmp.Diff(wantQuery, gotQuery))
	}
}

var (
	info = []byte(`{"server_id": "j230",
		"server_name": null,
		"expires": "2026-06-08 00:00:00",
		"expires_cytrus": null,
		"expires_storage": null,
		"param_ram": "1024",
		"param_disk": "10",
		"lastlog_panel": "2024-06-05 10:02:55",
		"mikrus_pro": "nie"
	}`)

	servers = []byte(`[
			{
				"server_id": "a133",
				"server_name": null,
				"expires": "2025-06-05 00:00:00",
				"param_ram": "1024",
				"param_disk": "10"
			},
			{
				"server_id": "j139",
				"server_name": null,
				"expires": "2026-06-08 00:00:00",
				"param_ram": "1024",
				"param_disk": "10"
			}
		]`)

	logs = []byte(`[
    {
        "id": "3752",
        "server_id": "j230",
        "task": "kluczssh",
        "when_created": "2024-06-05 10:05:34",
        "when_done": "2024-06-05 10:06:01",
        "output": "Wrzuciłem klucz SSH\n"
    },
    {
        "id": "3751",
        "server_id": "j230",
        "task": "restart",
        "when_created": "2024-06-05 09:57:54",
        "when_done": "2024-06-05 09:58:07",
        "output": "OK\n"
    },
    {
        "id": "3748",
        "server_id": "j230",
        "task": "upgrade",
        "when_created": "2024-06-05 08:59:28",
        "when_done": "2024-06-05 09:00:04",
        "output": "=== Aktualne parametry: 768 RAM / 10 DYSK\n2 / 20\nDodaje: +256MB RAM oraz +0GB dysku\nPo zmianie: 1024 MB / 10 GB\n[succes] GOTOWE!\n"
    }
]`)
)
