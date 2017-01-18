package kamemaru

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type List struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

func (k *kamemaru) route(conf config) error {
	// common middleware
	k.use()

	k.Echo.Static("/images", "images")
	k.Echo.POST("/register", k.register)
	k.Echo.POST("/login", k.login)

	// backend api
	api := k.Echo.Group("/api/v1")

	api.Use(middleware.JWT(k.JWTSecret))
	api.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	api.POST("/list", k.List)
	api.POST("/download", k.YoutubeDownload)

	return nil
}
