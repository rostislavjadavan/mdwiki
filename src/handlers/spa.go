package handlers

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"

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

func SPAHandler(e *echo.Echo) func(c echo.Context) error {
	return func(c echo.Context) error {
		content, err := fs.ReadFile(spaFS, "index.html")
		if err != nil {
			return echo.ErrNotFound
		}
		return c.HTMLBlob(http.StatusOK, content)
	}
}

func ViteDevProxyHandler(viteURL string) func(c echo.Context) error {
	target, err := url.Parse(viteURL)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	return func(c echo.Context) error {
		req := c.Request()
		req.URL.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.Host = target.Host
		proxy.ServeHTTP(c.Response(), req)
		return nil
	}
}
