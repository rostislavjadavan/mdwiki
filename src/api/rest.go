package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rostislavjadavan/mdwiki/src/config"
	"github.com/rostislavjadavan/mdwiki/src/search"
	"github.com/rostislavjadavan/mdwiki/src/storage"
)

type RestPage struct {
	Filename   string `json:"filename"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	RawContent string `json:"rawContent"`
}

type RestPageInfo struct {
	Filename string `json:"filename"`
	ModTime  string `json:"modTime"`
	Version  int64  `json:"version"`
}

func toRestPage(p *storage.Page) RestPage {
	return RestPage{
		Filename:   p.Filename,
		Name:       p.Name,
		Content:    p.Content,
		RawContent: p.RawContent,
	}
}

func toRestPageInfo(pi storage.PageInfo) RestPageInfo {
	return RestPageInfo{
		Filename: pi.Filename,
		ModTime:  pi.ModTime.Format("2006-01-02T15:04:05Z07:00"),
		Version:  pi.Version,
	}
}

func RestPageListHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		pages, err := s.PageList()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		result := make([]RestPageInfo, len(pages))
		for i, p := range pages {
			result[i] = toRestPageInfo(p)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func RestPageGetHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		filename := storage.FixPageExtension(c.Param("page"))
		page, err := s.Page(filename)
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, toRestPage(page))
	}
}

func RestPageCreateHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := new(FilenameRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}
		req.Filename = storage.FixPageExtension(req.Filename)

		if err := storage.ValidateFilename(req.Filename); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Message: err.Error()})
		}
		if s.PageExists(req.Filename) {
			return c.JSON(http.StatusConflict, ErrorResponse{Message: storage.SamePageExistsValidation})
		}

		page, err := s.PageCreate(req.Filename)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusCreated, toRestPage(page))
	}
}

func RestPageUpdateHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := new(UpdatePageRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}
		req.Filename = storage.FixPageExtension(req.Filename)

		if err := storage.ValidateFilename(req.Filename); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Message: err.Error()})
		}

		currentFilename := storage.FixPageExtension(c.Param("page"))
		if req.Filename != currentFilename && s.PageExists(req.Filename) {
			return c.JSON(http.StatusConflict, ErrorResponse{Message: storage.SamePageExistsValidation})
		}

		page, err := s.Page(currentFilename)
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}

		if err := s.PageContentUpdate(req.Content, page); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		if err := s.PageRename(req.Filename, page); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, toRestPage(page))
	}
}

func RestPageDeleteHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		filename := storage.FixPageExtension(c.Param("page"))
		page, err := s.Page(filename)
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}
		if err := s.PageDelete(page); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func RestSearchHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		query := c.QueryParam("query")
		result, err := search.Search(query, s)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, result)
	}
}

func RestTrashListHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		pages, err := s.TrashList()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		result := make([]RestPageInfo, len(pages))
		for i, p := range pages {
			result[i] = toRestPageInfo(p)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func RestTrashPageGetHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		filename := storage.FixPageExtension(c.Param("page"))
		page, err := s.TrashPage(filename)
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, toRestPage(page))
	}
}

func RestTrashRestoreHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		filename := storage.FixPageExtension(c.Param("page"))
		page, err := s.TrashPage(filename)
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}
		if err := s.TrashRestore(page); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"filename": page.Filename})
	}
}

func RestTrashEmptyHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		if err := s.TrashEmpty(); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func RestVersionsListHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		filename := storage.FixPageExtension(c.Param("page"))
		versions, err := s.VersionsList(filename)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		result := make([]RestPageInfo, len(versions))
		for i, v := range versions {
			result[i] = toRestPageInfo(v)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func RestVersionGetHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		ver := c.Param("ver")
		page, err := s.VersionPage(ver)
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, toRestPage(page))
	}
}

func RestSettingsHandler(cfg *config.AppConfig) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"theme": cfg.Theme})
	}
}

func RestVersionRestoreHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		ver := c.Param("ver")
		page, err := s.VersionPage(ver)
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}
		if err := s.VersionRestore(page); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"name": page.Name})
	}
}
