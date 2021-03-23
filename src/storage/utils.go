package storage

import (
	"errors"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

var defaultFilePermission fs.FileMode = 0755
var defaultDirectoryPermission fs.FileMode = 0755

var renderer markdown.Renderer

var filenameValidationRegexp *regexp.Regexp

var FilenameEmptyValidation string = "Filename cannot be empty"
var InvalidFilenameValidation string = "Invalid filename, valid examples: wiki_page_1.md, flowers-and-animals.md, page106.md"
var SamePageExistsValidation string = "Page with same name already exists"

func init() {
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer = html.NewRenderer(opts)

	filenameValidationRegexp, _ = regexp.Compile("^[\\w-\\\\.]+[^\\W]{1}$")
}

func ToMarkdown(content []byte) string {
	contentString := strings.NewReplacer("\r\n", "\n").Replace(string(content))
	return string(markdown.ToHTML([]byte(contentString), nil, renderer))
}

func FixPageExtension(uri string) string {
	page := strings.TrimSpace(uri)
	if page != "" && !strings.HasSuffix(page, ".md") {
		page = page + ".md"
	}
	return page
}

func ValidateFilename(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New(FilenameEmptyValidation)
	}
	if !filenameValidationRegexp.MatchString(filename) {
		return errors.New(InvalidFilenameValidation)
	}
	return nil
}

func fsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
