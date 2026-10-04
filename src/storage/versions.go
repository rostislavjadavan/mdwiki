package storage

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Storage) savePageVersion(page *Page) error {
	filename := page.Filename + "__" + strconv.FormatInt(time.Now().Unix(), 10)
	err := ioutil.WriteFile(path.Join(s.config.Storage, "versions", filename), []byte(page.RawContent), defaultFilePermission)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) renamePageVersions(oldFilename string, newFilename string) error {
	files, err := filepath.Glob(path.Join(s.config.Storage, "versions", oldFilename+"__*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Rename(file, strings.Replace(file, oldFilename, newFilename, 1))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) VersionsList(filename string) ([]PageInfo, error) {
	dir, err := os.Open(path.Join(s.config.Storage, "versions"))
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
		if !file.IsDir() && ValidateFilename(file.Name()) == nil && strings.HasPrefix(file.Name(), filename+"__") {
			_, version := parsePageFilename(file.Name())
			list = append(list, PageInfo{
				Filename: file.Name(),
				ModTime:  file.ModTime(),
				Version:  version,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Version > list[j].Version
	})

	return list, nil
}

func (s *Storage) VersionPage(filename string) (*Page, error) {
	content, err := ioutil.ReadFile(path.Join(s.config.Storage, "versions", filename))
	if err != nil {
		return nil, err
	}
	return newPage(filename, content), nil
}

func (s *Storage) VersionRestore(page *Page) error {
	currentPage, err := s.Page(page.Name)
	if err != nil {
		return err
	}
	return s.PageContentUpdate(page.RawContent, currentPage)
}
