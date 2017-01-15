package kamemaru

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/lestrrat/go-server-starter/listener"
	"github.com/uber-go/zap"
)

var (
	DeployMode string
	LogPath    string
)

type kamemaru struct {
	Echo   *echo.Echo
	Logger zap.Logger
}

func New() *kamemaru {
	kame := &kamemaru{Echo: echo.New()}
	return kame.setup()
}

func (k *kamemaru) Run() int {
	k.Route()
	if err := k.RunServer(); err != nil {
		k.Logger.Error("Failed to run server", zap.String("reason", err.Error()))
		return 1
	}
	return 0
}

func (k *kamemaru) RunServer() error {
	var l net.Listener

	port := os.Getenv("SERVER_STARTER_PORT")
	if port != "" {
		listeners, err := listener.ListenAll()
		if err != nil {
			return err
		}

		if len(listeners) > 0 {
			l = listeners[0]
		}
	}

	if l == nil {
		var err error
		port = ":8080"
		l, err = net.Listen("tcp", port)
		if err != nil {
			return err
		}
	}

	k.Logger.Info("Graceful start kamemaru...", zap.String("port", port))

	return serve(k.Echo.Server, l)
}

func (k *kamemaru) setup() *kamemaru {
	switch DeployMode {
	case "develop":
		k.setlogger(os.Stderr)
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

		k.setlogger(zap.AddSync(f))
	default:
		log.Fatal("kamemaru.Deploymode was not set")
	}

	k.use()

	return k
}

func (k *kamemaru) setlogger(Out zap.WriteSyncer) {
	k.Logger = zap.New(
		zap.NewJSONEncoder(JSTFormatter("time")),
		zap.AddCaller(), // Add line number option
		zap.Output(Out),
	)
}

func serve(server *http.Server, l net.Listener) error {
	return server.Serve(l)
}
