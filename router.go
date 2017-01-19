package kamemaru

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (k *kamemaru) route(conf config) error {
	// common middleware
	k.use()

	k.Echo.Static("/images", "images")
	k.Echo.POST("/register", k.register)
	k.Echo.POST("/login", k.login)

	// backend api
	k.Echo.GET("/api/fetch", k.imgfetch)

	api := k.Echo.Group("/api/v1")
	api.Use(middleware.JWT(k.JWTSecret))
	api.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	api.POST("/upload", k.Upload)

	return nil
}
