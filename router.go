package kamemaru

import "github.com/labstack/echo/middleware"

type List struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

func (k *kamemaru) route(conf config) error {
	// common middleware
	k.use()

	k.Echo.POST("/register", k.register)
	k.Echo.POST("/login", k.login)

	// backend api
	api := k.Echo.Group("/api/v1")

	api.Use(middleware.JWT(k.JWTSecret))
	api.POST("/api/v1/list", k.List)
	api.POST("/api/v1/download", k.YoutubeDownload)

	return nil
}
