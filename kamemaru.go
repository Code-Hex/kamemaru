package kamemaru

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/uber-go/zap"
)

var (
	DeployMode string
	LogPath    string
	Port       string
)

type kamemaru struct {
	Echo   *echo.Echo
	Logger zap.Logger
}

func New() *kamemaru {
	kame := &kamemaru{Echo: echo.New()}
	kame.SetDeployMode()

	return kame
}

func (k *kamemaru) Run() int {
	k.Route()
	return 0
}

func (k *kamemaru) SetDeployMode() {
	switch DeployMode {
	case "develop":
		k.SetLogger(os.Stderr)
	case "staging":
		err := os.MkdirAll(LogPath, 0755)
		if err != nil {
			log.Fatal(err.Error())
		}

		f, err := rotatelogs.New(
			filepath.Join(LogPath, "access_log.%Y%m%d%H%ls M"),
			rotatelogs.WithLinkName(filepath.Join(LogPath, "access_log")),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil {
			log.Fatalf("failed to create rotatelogs: %s", err)
		}

		defer f.Close()

		k.SetLogger(zap.AddSync(f))
	default:
		log.Fatal("kamemaru.Deploymode was not set")
	}

	k.Logger.Info("Graceful start kamemaru...", zap.String("Port", Port))
}

func (k *kamemaru) SetLogger(Out zap.WriteSyncer) {
	k.Logger = zap.New(
		zap.NewJSONEncoder(JSTFormatter("time")),
		zap.AddCaller(), // Add line number option
		zap.Output(Out),
	)

	k.Echo.Use(k.loghandler())
}

func (k *kamemaru) loghandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if err = next(c); err != nil {
				c.Error(err)
			}

			w, r := c.Response(), c.Request()
			k.Logger.Info(
				"Detect access",
				zap.Int("statuscode", w.Status),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("useragent", r.UserAgent()),
				zap.String("remote_ip", r.RemoteAddr),
			)
			return nil
		}
	}
}

func JSTFormatter(key string) zap.TimeFormatter {
	return zap.TimeFormatter(func(t time.Time) zap.Field {
		jst := time.FixedZone("Asia/Tokyo", 9*3600)
		return zap.String(key, t.In(jst).Format(time.ANSIC))
	})
}
