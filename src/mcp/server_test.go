package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rostislavjadavan/mdwiki/src/config"
	"github.com/rostislavjadavan/mdwiki/src/storage"
)

func testServer(t *testing.T, token string) (*storage.Storage, string, *httptest.Server) {
	t.Helper()
	dir := t.TempDir()
	s, err := storage.CreateStorage(&config.AppConfig{Storage: dir})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(NewHandler(s, &sync.Mutex{}, token))
	t.Cleanup(server.Close)
	return s, dir, server
}

type tokenTransport struct{ token string }

func (tr tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+tr.token)
	return http.DefaultTransport.RoundTrip(req)
}

func TestToolsOverHTTP(t *testing.T) {
	_, dir, server := testServer(t, "secret")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "1"}, nil)
	session, err := client.Connect(ctx, &sdk.StreamableClientTransport{
		Endpoint: server.URL, DisableStandaloneSSE: true,
		HTTPClient: &http.Client{Transport: tokenTransport{token: "secret"}},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 13 {
		t.Fatalf("got %d tools", len(tools.Tools))
	}
	for _, tool := range tools.Tools {
		if tool.InputSchema == nil || tool.Annotations == nil {
			t.Fatalf("missing metadata for %s", tool.Name)
		}
	}
	call := func(name string, args map[string]any, wantError bool) json.RawMessage {
		t.Helper()
		result, err := session.CallTool(ctx, &sdk.CallToolParams{Name: name, Arguments: args})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if result.IsError != wantError {
			t.Fatalf("%s: error=%v, content=%v", name, result.IsError, result.Content)
		}
		if wantError {
			return nil
		}
		if len(result.Content) != 1 {
			t.Fatalf("%s: unexpected content", name)
		}
		text, ok := result.Content[0].(*sdk.TextContent)
		if !ok || !json.Valid([]byte(text.Text)) {
			t.Fatalf("%s: invalid JSON content", name)
		}
		return json.RawMessage(text.Text)
	}
	page := func(name string) map[string]any { return map[string]any{"filename": name} }
	content := func(name, text string) map[string]any { return map[string]any{"filename": name, "content": text} }
	call("list_pages", map[string]any{}, false)
	call("create_page", content("test", "# First\nUnique phrase"), false)
	call("create_page", content("test", "overwrite"), true)
	call("read_page", page("test.md"), false)
	call("search_pages", map[string]any{"query": "Unique phrase"}, false)
	call("update_page", content("test", "# Second"), false)
	call("list_versions", page("test"), false)
	// Seed a stable historical version independent of second-resolution saves.
	version := "test.md__123"
	if err := os.WriteFile(filepath.Join(dir, "versions", version), []byte("# Historical"), 0600); err != nil {
		t.Fatal(err)
	}
	call("read_version", map[string]any{"version": version}, false)
	data := call("restore_version", map[string]any{"version": version}, false)
	var restored pageOutput
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Content != "# Historical" {
		t.Fatalf("restored %q", restored.Content)
	}
	call("rename_page", map[string]any{"filename": "test", "newFilename": "home"}, true)
	call("rename_page", map[string]any{"filename": "test", "newFilename": "renamed"}, false)
	call("read_page", page("test"), true)
	call("delete_page", page("renamed"), false)
	call("list_trash", map[string]any{}, false)
	call("read_trash_page", page("renamed"), false)
	call("create_page", content("renamed", "conflicting active page"), false)
	call("restore_page", page("renamed"), true)
	call("delete_page", page("renamed"), true)
	// Remove only the conflicting fixture so the original page can be restored.
	if err := os.Remove(filepath.Join(dir, "pages", "renamed.md")); err != nil {
		t.Fatal(err)
	}
	call("restore_page", page("renamed"), false)
	call("read_page", page("renamed"), false)
	for _, tool := range []string{"read_page", "create_page", "update_page", "delete_page", "restore_page", "list_versions", "read_trash_page"} {
		args := content("../outside", "bad")
		call(tool, args, true)
	}
	call("rename_page", map[string]any{"filename": "renamed", "newFilename": "../outside"}, true)
	call("read_version", map[string]any{"version": "../outside.md__123"}, true)
	call("restore_version", map[string]any{"version": "../outside.md__123"}, true)
	call("update_page", page("renamed"), true) // Missing required content.
}

func TestHTTPAccessChecks(t *testing.T) {
	_, _, server := testServer(t, "secret")
	for _, tc := range []struct {
		name, token, origin string
		status              int
	}{
		{"missing token", "", "", http.StatusUnauthorized},
		{"wrong token", "wrong", "", http.StatusUnauthorized},
		{"foreign origin", "secret", "https://evil.example", http.StatusForbidden},
		{"null origin", "secret", "null", http.StatusForbidden},
		{"same origin", "secret", server.URL, http.StatusMethodNotAllowed},
		{"nonbrowser", "secret", "", http.StatusMethodNotAllowed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "text/event-stream")
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.status {
				t.Fatalf("got %d, want %d", resp.StatusCode, tc.status)
			}
		})
	}
}

func TestFilenameValidation(t *testing.T) {
	for _, name := range []string{"", "../home", "a/b", "a\\b", "..", "a\x00b"} {
		if _, err := filename(name); err == nil {
			t.Errorf("accepted filename %q", name)
		}
	}
	for _, name := range []string{"../home.md__123", "home__123", "home.md__bad", "home.md", "home.md__123/abc"} {
		if err := versionName(name); err == nil {
			t.Errorf("accepted version %q", name)
		}
	}
}
