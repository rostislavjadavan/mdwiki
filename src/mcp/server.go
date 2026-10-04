// Package mcp exposes the wiki through the Model Context Protocol.
package mcp

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rostislavjadavan/mdwiki/src/search"
	"github.com/rostislavjadavan/mdwiki/src/storage"
)

type emptyInput struct{}
type pageInput struct {
	Filename string `json:"filename" jsonschema:"Page filename, with or without .md"`
}
type contentInput struct {
	Filename string `json:"filename" jsonschema:"Page filename, with or without .md"`
	Content  string `json:"content" jsonschema:"Complete replacement Markdown content"`
}
type renameInput struct {
	Filename    string `json:"filename"`
	NewFilename string `json:"newFilename"`
}
type versionInput struct {
	Version string `json:"version" jsonschema:"Exact version filename returned by list_versions"`
}
type searchInput struct {
	Query string `json:"query"`
}
type pageOutput struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func pageResult(p *storage.Page) pageOutput {
	return pageOutput{Filename: p.Filename, Content: p.RawContent}
}

func filename(name string) (string, error) {
	if strings.TrimSpace(name) == "." || strings.TrimSpace(name) == ".." {
		return "", fmt.Errorf("invalid filename")
	}
	name = storage.FixPageExtension(name)
	if strings.ContainsAny(name, "/\\") {
		return "", fmt.Errorf("invalid filename")
	}
	return name, storage.ValidateFilename(name)
}

func versionName(name string) error {
	i := strings.LastIndex(name, "__")
	if i < 0 {
		return fmt.Errorf("invalid version filename")
	}
	base, err := filename(name[:i])
	if err != nil || base != name[:i] {
		return fmt.Errorf("invalid version filename")
	}
	_, err = strconv.ParseInt(name[i+2:], 10, 64)
	return err
}

// addTool serializes storage access with REST requests, and lets the SDK
// validate the inferred input schema and report errors as MCP tool errors.
func addTool[I any](server *mcp.Server, mu *sync.Mutex, name, description string, readOnly bool, fn func(I) (any, error)) {
	destructive := !readOnly
	mcp.AddTool(server, &mcp.Tool{
		Name: name, Description: description,
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: readOnly, DestructiveHint: &destructive},
	}, func(ctx context.Context, req *mcp.CallToolRequest, input I) (*mcp.CallToolResult, any, error) {
		mu.Lock()
		defer mu.Unlock()
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}
		result, err := fn(input)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.Marshal(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
	})
}

// NewHandler creates a stateless Streamable HTTP endpoint. The same mutex
// must also guard REST storage operations.
func NewHandler(s *storage.Storage, mu *sync.Mutex, token string) http.Handler {
	server := mcp.NewServer(&mcp.Implementation{Name: "mdwiki", Version: "1.0.0"}, nil)
	addTool(server, mu, "list_pages", "List active wiki pages.", true, func(emptyInput) (any, error) { return s.PageList() })
	addTool(server, mu, "read_page", "Read a wiki page as Markdown.", true, func(in pageInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		p, err := s.Page(name)
		if err != nil {
			return nil, err
		}
		return pageResult(p), nil
	})
	addTool(server, mu, "search_pages", "Fuzzy search page filenames and content. Previews contain HTML highlighting.", true, func(in searchInput) (any, error) { return search.Search(in.Query, s) })
	addTool(server, mu, "create_page", "Create a new Markdown page. Refuses to overwrite an existing page.", false, func(in contentInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		if s.PageExists(name) {
			return nil, fmt.Errorf("page already exists")
		}
		p, err := s.PageCreate(name)
		if err != nil {
			return nil, err
		}
		if err := s.PageContentUpdate(in.Content, p); err != nil {
			return nil, err
		}
		return pageResult(p), nil
	})
	addTool(server, mu, "update_page", "Replace a page's complete Markdown content, preserving its previous content as a version.", false, func(in contentInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		p, err := s.Page(name)
		if err != nil {
			return nil, err
		}
		if err := s.PageContentUpdate(in.Content, p); err != nil {
			return nil, err
		}
		return pageResult(p), nil
	})
	addTool(server, mu, "rename_page", "Rename a page and its versions. Refuses to overwrite an existing page.", false, func(in renameInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		newName, err := filename(in.NewFilename)
		if err != nil {
			return nil, err
		}
		p, err := s.Page(name)
		if err != nil {
			return nil, err
		}
		if name != newName && s.PageExists(newName) {
			return nil, fmt.Errorf("target page already exists")
		}
		if err := s.PageRename(newName, p); err != nil {
			return nil, err
		}
		return pageResult(p), nil
	})
	addTool(server, mu, "delete_page", "Move a page to trash. Refuses to overwrite a trashed page with the same name.", false, func(in pageInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		trash, err := s.TrashList()
		if err != nil {
			return nil, err
		}
		for _, p := range trash {
			if p.Filename == name {
				return nil, fmt.Errorf("page already exists in trash")
			}
		}
		p, err := s.Page(name)
		if err != nil {
			return nil, err
		}
		return map[string]string{"filename": name}, s.PageDelete(p)
	})
	addTool(server, mu, "list_trash", "List deleted pages that can be restored.", true, func(emptyInput) (any, error) { return s.TrashList() })
	addTool(server, mu, "read_trash_page", "Read a deleted page as Markdown.", true, func(in pageInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		p, err := s.TrashPage(name)
		if err != nil {
			return nil, err
		}
		return pageResult(p), nil
	})
	addTool(server, mu, "restore_page", "Restore a page from trash. Refuses to overwrite an active page.", false, func(in pageInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		if s.PageExists(name) {
			return nil, fmt.Errorf("active page already exists")
		}
		p, err := s.TrashPage(name)
		if err != nil {
			return nil, err
		}
		if err := s.TrashRestore(p); err != nil {
			return nil, err
		}
		return pageResult(p), nil
	})
	addTool(server, mu, "list_versions", "List a page's saved versions. Use the returned Filename for read_version or restore_version.", true, func(in pageInput) (any, error) {
		name, err := filename(in.Filename)
		if err != nil {
			return nil, err
		}
		return s.VersionsList(name)
	})
	addTool(server, mu, "read_version", "Read a saved version as Markdown.", true, func(in versionInput) (any, error) {
		if err := versionName(in.Version); err != nil {
			return nil, err
		}
		p, err := s.VersionPage(in.Version)
		if err != nil {
			return nil, err
		}
		return pageResult(p), nil
	})
	addTool(server, mu, "restore_version", "Restore a saved version to its active page, preserving the current content as a version. The active page must exist.", false, func(in versionInput) (any, error) {
		if err := versionName(in.Version); err != nil {
			return nil, err
		}
		p, err := s.VersionPage(in.Version)
		if err != nil {
			return nil, err
		}
		if err := s.VersionRestore(p); err != nil {
			return nil, err
		}
		current, err := s.Page(p.Name)
		if err != nil {
			return nil, err
		}
		return pageResult(current), nil
	})
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			u, err := url.Parse(origin)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host != r.Host || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
				http.Error(w, "forbidden origin", http.StatusForbidden)
				return
			}
		}
		if token != "" && subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+token)) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<20)
		handler.ServeHTTP(w, r)
	})
}
