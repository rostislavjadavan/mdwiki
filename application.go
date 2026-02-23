package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rostislavjadavan/mdwiki/src/api"
	"github.com/rostislavjadavan/mdwiki/src/config"
	"github.com/rostislavjadavan/mdwiki/src/handlers"
	"github.com/rostislavjadavan/mdwiki/src/storage"
	"github.com/rostislavjadavan/mdwiki/src/ui"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	s, err := storage.CreateStorage(cfg)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	// Static files
	e.GET("/static/style.css", handlers.StaticHandler(ui.CssStyle, handlers.MimeCss, e))
	e.GET("/static/script.js", handlers.StaticHandler(ui.JavascriptScript, handlers.MimeJavascript, e))
	e.GET("/static/favicon.png", handlers.StaticHandler(ui.ImageFaviconPng, handlers.MimePng, e))

	// RPC like API
	e.POST("/api/page.create", api.PageCreateHandler(e, s))
	e.POST("/api/page.update/:page", api.PageUpdateHandler(e, s))
	e.POST("/api/page.delete", api.PageDeleteHandler(e, s))
	e.POST("/api/trash.empty", api.TrashEmptyHandler(e, s))
	e.POST("/api/trash.restore", api.TrashRestoreHandler(e, s))
	e.POST("/api/version.restore", api.VersionRestoreHandler(e, s))

	// REST API
	e.GET("/api/pages", api.RestPageListHandler(e, s))
	e.GET("/api/pages/:page", api.RestPageGetHandler(e, s))
	e.POST("/api/pages", api.RestPageCreateHandler(e, s))
	e.PUT("/api/pages/:page", api.RestPageUpdateHandler(e, s))
	e.DELETE("/api/pages/:page", api.RestPageDeleteHandler(e, s))
	e.GET("/api/search", api.RestSearchHandler(e, s))
	e.GET("/api/trash", api.RestTrashListHandler(e, s))
	e.GET("/api/trash/:page", api.RestTrashPageGetHandler(e, s))
	e.POST("/api/trash/:page/restore", api.RestTrashRestoreHandler(e, s))
	e.DELETE("/api/trash", api.RestTrashEmptyHandler(e, s))
	e.GET("/api/pages/:page/versions", api.RestVersionsListHandler(e, s))
	e.GET("/api/versions/:ver", api.RestVersionGetHandler(e, s))
	e.POST("/api/versions/:ver/restore", api.RestVersionRestoreHandler(e, s))

	// UI
	e.GET("/search", handlers.SearchHandler(e, s))
	e.GET("/trash", handlers.TrashHandler(e, s))
	e.GET("/trash/:page", handlers.TrashPageHandler(e, s))
	e.GET("/list", handlers.ListHandler(e, s))
	e.GET("/create", handlers.CreateHandler(e, s))
	e.GET("/edit/:page", handlers.EditHandler(e, s))
	e.GET("/", handlers.PageHandler(e, s))
	e.GET("/:page/version", handlers.PageVersionsHandler(e, s))
	e.GET("/:page/version/:ver", handlers.PageVersionHandler(e, s))
	e.GET("/:page", handlers.PageHandler(e, s))

	// SPA
	e.GET("/static/app/*", handlers.SPAStaticHandler(e))
	e.GET("/*", handlers.SPAHandler(e))

	e.Logger.Fatal(e.Start(cfg.Host + ":" + cfg.Port))
}
