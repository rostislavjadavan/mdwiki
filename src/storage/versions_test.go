package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rostislavjadavan/mdwiki/src/config"
)

func TestVersionsListMatchesExactPage(t *testing.T) {
	dir := t.TempDir()
	s, err := CreateStorage(&config.AppConfig{Storage: dir})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"test.md__123", "other-test.md__124", "test.md-extra.md__125"} {
		if err := os.WriteFile(filepath.Join(dir, "versions", name), []byte("version"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	versions, err := s.VersionsList("test.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(versions) != 1 || versions[0].Filename != "test.md__123" {
		t.Fatalf("unexpected versions: %+v", versions)
	}
}

func TestVersionRestoreReturnsSaveError(t *testing.T) {
	dir := t.TempDir()
	s, err := CreateStorage(&config.AppConfig{Storage: dir})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.PageCreate("test.md"); err != nil {
		t.Fatal(err)
	}
	// Redirect the versions directory to a regular file, forcing snapshot failure.
	if err := os.Rename(filepath.Join(dir, "versions"), filepath.Join(dir, "saved-versions")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "versions"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := s.VersionRestore(newPage("test.md__123", []byte("historical"))); err == nil {
		t.Fatal("restore hid the snapshot failure")
	}
	p, err := s.Page("test.md")
	if err != nil {
		t.Fatal(err)
	}
	if p.RawContent != "# test.md" {
		t.Fatal("failed restore changed the page")
	}
}
