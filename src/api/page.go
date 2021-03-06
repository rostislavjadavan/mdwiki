package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rostislavjadavan/mdwiki/src/storage"
	"net/http"
)

func PageCreateHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := new(FilenameRequest)
		err := c.Bind(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}
		req.Filename = storage.FixPageExtension(req.Filename)

		e.Logger.Debugf("api create '%s'", req.Filename)

		// Filename validation
		err = storage.ValidateFilename(req.Filename)
		if err != nil {
			return c.JSON(http.StatusExpectationFailed, ErrorResponse{Message: err.Error()})
		}

		if s.PageExists(req.Filename) {
			return c.JSON(http.StatusExpectationFailed, ErrorResponse{Message: storage.SamePageExistsValidation})
		}
		// Filename validation end

		e.Logger.Debug("creating new page " + req.Filename)
		page, err := s.PageCreate(req.Filename)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, RedirectResponse{Redirect: "/edit/" + page.Filename})
	}
}

func PageUpdateHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		e.Logger.Debugf("api update '%s'", c.Param("page"))

		req := new(UpdatePageRequest)
		err := c.Bind(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}
		req.Filename = storage.FixPageExtension(req.Filename)

		// Filename validation
		err = storage.ValidateFilename(req.Filename)
		if err != nil {
			return c.JSON(http.StatusExpectationFailed, ErrorResponse{Message: err.Error()})
		}

		if req.Filename != c.Param("page") && s.PageExists(req.Filename) {
			return c.JSON(http.StatusExpectationFailed, ErrorResponse{Message: storage.SamePageExistsValidation})
		}
		// Filename validation end

		page, err := s.Page(c.Param("page"))
		if err != nil {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}

		err = s.PageContentUpdate(req.Content, page)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}

		err = s.PageRename(req.Filename, page)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, RedirectResponse{Redirect: "/" + page.Filename})
	}
}

func PageDeleteHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := new(FilenameRequest)
		err := c.Bind(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}
		req.Filename = storage.FixPageExtension(req.Filename)

		e.Logger.Debugf("api delete '%s'", req.Filename)

		page, err := s.Page(storage.FixPageExtension(req.Filename))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		err = s.PageDelete(page)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, RedirectResponse{Redirect: "/list"})
	}
}
