package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/rostislavjadavan/mdwiki/src/storage"
	"github.com/rostislavjadavan/mdwiki/src/ui"
	"net/http"
)

var HomePage string = "home.md"

func PageHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		e.Logger.Debug("page /" + c.Param("page") + " requested")

		pageUri := c.Param("page")
		if pageUri == HomePage {
			return c.Redirect(http.StatusPermanentRedirect, "/")
		}
		if pageUri == "" {
			pageUri = HomePage
		}
		if pageUri != storage.FixPageExtension(pageUri) {
			return c.Redirect(http.StatusPermanentRedirect, "/"+storage.FixPageExtension(pageUri))
		}

		page, err := s.LoadPage(pageUri)
		if err != nil {
			return notFoundPage(err, e, c)
		}

		tpl, err := ui.Render(ui.TemplatePage, page)
		if err != nil {
			return errorPage(err, e, c)
		}

		return c.HTML(http.StatusOK, tpl)
	}
}

func ListHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		e.Logger.Debug("page list")

		list, err := s.ListPages()
		if err != nil {
			return errorPage(err, e, c)
		}

		tpl, err := ui.Render(ui.TemplateList, map[string]interface{}{
			"Pages": list,
		})
		if err != nil {
			return errorPage(err, e, c)
		}

		return c.HTML(http.StatusOK, tpl)
	}
}

func CreateHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		tpl, err := ui.Render(ui.TemplateCreate, nil)
		if err != nil {
			return errorPage(err, e, c)
		}

		return c.HTML(http.StatusOK, tpl)
	}
}

func EditHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		e.Logger.Debug("edit /" + c.Param("page"))

		page, err := s.LoadPage(storage.FixPageExtension(c.Param("page")))
		if err != nil {
			return notFoundPage(err, e, c)
		}

		tpl, err := ui.Render(ui.TemplateEdit, page)
		if err != nil {
			return errorPage(err, e, c)
		}

		return c.HTML(http.StatusOK, tpl)
	}
}

func DeleteHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		e.Logger.Debug("delete /" + c.Param("page"))

		page, err := s.LoadPage(storage.FixPageExtension(c.Param("page")))
		if err != nil {
			return notFoundPage(err, e, c)
		}

		tpl, err := ui.Render(ui.TemplateDelete, page)
		if err != nil {
			return errorPage(err, e, c)
		}

		return c.HTML(http.StatusOK, tpl)
	}
}

func DoDeleteHandler(e *echo.Echo, s *storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		page, err := s.LoadPage(storage.FixPageExtension(c.Param("page")))
		if err != nil {
			return notFoundPage(err, e, c)
		}
		err = s.DeletePage(page)
		if err != nil {
			return errorPage(err, e, c)
		}
		return c.Redirect(http.StatusFound, "/list")
	}
}
