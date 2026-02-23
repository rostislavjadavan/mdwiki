package handlers

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	mdwikifrontend "github.com/rostislavjadavan/mdwiki/frontend"
)

var spaFS fs.FS

func init() {
	var err error
	spaFS, err = fs.Sub(mdwikifrontend.Dist, "dist")
	if err != nil {
		panic(err)
	}
}

func SPAStaticHandler(e *echo.Echo) func(c echo.Context) error {
	fileServer := http.FileServer(http.FS(spaFS))
	return func(c echo.Context) error {
		c.Request().URL.Path = c.Param("*")
		fileServer.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func SPAFaviconHandler(e *echo.Echo) func(c echo.Context) error {
	return func(c echo.Context) error {
		content, err := fs.ReadFile(spaFS, "favicon.ico")
		if err != nil {
			// Fall back to vite.svg or return 404
			content, err = fs.ReadFile(spaFS, "vite.svg")
			if err != nil {
				return echo.ErrNotFound
			}
			c.Response().Header().Set("Content-Type", MimeSvg)
			return c.Blob(http.StatusOK, MimeSvg, content)
		}
		c.Response().Header().Set("Content-Type", MimeIco)
		return c.Blob(http.StatusOK, MimeIco, content)
	}
}

func SPAHandler(e *echo.Echo) func(c echo.Context) error {
	return func(c echo.Context) error {
		content, err := fs.ReadFile(spaFS, "index.html")
		if err != nil {
			return echo.ErrNotFound
		}
		return c.HTMLBlob(http.StatusOK, content)
	}
}
