package storage

import (
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

func (s *Storage) PageCreate(filename string) (*Page, error) {
	if !strings.HasSuffix(filename, ".md") {
		filename = filename + ".md"
	}
	err := ioutil.WriteFile(path.Join(s.config.Storage, "pages", filename), []byte("# "+filename), defaultFilePermission)
	if err != nil {
		return nil, err
	}
	return s.Page(filename)
}

func (s *Storage) Page(filename string) (*Page, error) {
	content, err := ioutil.ReadFile(path.Join(s.config.Storage, "pages", filename))
	if err != nil {
		return nil, err
	}
	return newPage(filename, content), nil
}

func (s *Storage) PageRawContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(path.Join(s.config.Storage, "pages", filename))
	if err != nil {
		return "", err
	}
	return string(content[:]), nil
}

func (s *Storage) PageDelete(page *Page) error {
	return os.Rename(path.Join(s.config.Storage, "pages", page.Filename), path.Join(s.config.Storage, "trash", page.Filename))
}

func (s *Storage) TrashRestore(page *Page) error {
	return os.Rename(path.Join(s.config.Storage, "trash", page.Filename), path.Join(s.config.Storage, "pages", page.Filename))
}

func (s *Storage) PageContentUpdate(content string, page *Page) error {
	if content != page.RawContent {
		err := s.savePageVersion(page)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path.Join(s.config.Storage, "pages", page.Filename), []byte(content), defaultFilePermission)
		if err != nil {
			return err
		}

		page.RawContent = content
		page.Content = ToMarkdown([]byte(content))
	}
	return nil
}

func (s *Storage) PageRename(newFilename string, page *Page) error {
	if newFilename != page.Filename {
		err := s.renamePageVersions(page.Filename, newFilename)
		if err != nil {
			return err
		}

		err = os.Rename(path.Join(s.config.Storage, "pages", page.Filename), path.Join(s.config.Storage, "pages", newFilename))
		if err != nil {
			return err
		}

		page.Filename = newFilename
	}
	return nil
}

type PageInfo struct {
	Filename string
	ModTime  time.Time
	Version  int64
}

func (s *Storage) PageList() ([]PageInfo, error) {
	dir, err := os.Open(path.Join(s.config.Storage, "pages"))
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	list := make([]PageInfo, 0)
	for _, file := range files {
		if !file.IsDir() && ValidateFilename(file.Name()) == nil {
			list = append(list, PageInfo{
				Filename: file.Name(),
				ModTime:  file.ModTime(),
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Filename < list[j].Filename
	})

	return list, nil
}

func (s *Storage) PageExists(filename string) bool {
	_, err := s.Page(filename)
	if err == nil {
		return true
	}
	return false
}
