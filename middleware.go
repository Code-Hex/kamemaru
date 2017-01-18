package kamemaru

import (
	"fmt"
	"net/http"
	"time"

	static "github.com/Code-Hex/echo-static"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/uber-go/zap"
)

func (k *kamemaru) use() {
	k.Echo.HTTPErrorHandler = k.ErrorHandler

	k.Echo.Use(k.LogHandler())
	k.Echo.Use(static.ServeRoot("/static", NewAssets("assets")))
	k.Echo.Use(middleware.Recover())
}

func (k *kamemaru) LogHandler() echo.MiddlewareFunc {
	return func(before echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := before(c)
			stop := time.Now()

			w, r := c.Response(), c.Request()
			k.Logger.Info(
				"Detected access",
				zap.String("status", fmt.Sprintf("%d: %s", w.Status, http.StatusText(w.Status))),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("useragent", r.UserAgent()),
				zap.String("remote_ip", r.RemoteAddr),
				zap.Int64("latency", stop.Sub(start).Nanoseconds()/int64(time.Microsecond)),
			)
			return err
		}
	}
}

func (k *kamemaru) ErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // Issue #608
			if err := c.NoContent(code); err != nil {
				goto ERROR
			}
		} else {
			if err := c.JSON(code, msg); err != nil {
				goto ERROR
			}
		}
	}
ERROR:
	k.Logger.Error("Error", zap.String("reason", err.Error()))
}

func JSTFormatter(key string) zap.TimeFormatter {
	return zap.TimeFormatter(func(t time.Time) zap.Field {
		jst := time.FixedZone("Asia/Tokyo", 9*3600)
		return zap.String(key, t.In(jst).Format(time.ANSIC))
	})
}

func NewAssets(root string) *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    root,
	}
}
