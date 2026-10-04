package storage

import (
	"io/ioutil"
	"os"
	"path"
	"sort"
)

func (s *Storage) TrashPage(filename string) (*Page, error) {
	content, err := ioutil.ReadFile(path.Join(s.config.Storage, "trash", filename))
	if err != nil {
		return nil, err
	}
	return newPage(filename, content), nil
}

func (s *Storage) TrashList() ([]PageInfo, error) {
	dir, err := os.Open(path.Join(s.config.Storage, "trash"))
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

func (s *Storage) TrashEmpty() error {
	pages, err := s.TrashList()
	if err != nil {
		return err
	}
	for _, p := range pages {
		// Remove files from trash
		err := os.Remove(path.Join(s.config.Storage, "trash", p.Filename))
		if err != nil {
			return err
		}

		// Remove versions of files in trash
		versions, err := s.VersionsList(p.Filename)
		if err != nil {
			return err
		}
		for _, v := range versions {
			err := os.Remove(path.Join(s.config.Storage, "versions", v.Filename))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
