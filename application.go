package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rostislavjadavan/mdwiki/src/api"
	"github.com/rostislavjadavan/mdwiki/src/config"
	"github.com/rostislavjadavan/mdwiki/src/handlers"
	"github.com/rostislavjadavan/mdwiki/src/storage"
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
	e.GET("/api/settings", api.RestSettingsHandler(cfg))

	// SPA static assets and catch-all
	e.GET("/static/app/*", handlers.SPAStaticHandler(e))
	e.Any("/*", handlers.SPAHandler(e))

	e.Logger.Fatal(e.Start(cfg.Host + ":" + cfg.Port))
}
