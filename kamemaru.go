package kamemaru

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/uber-go/zap"
)

type Env struct {
	DeployMode string
	Logger     *zap.Logger
	Port       int
}

type kamemaru struct {
	Echo *echo.Echo
	Env
}

func New() *kamemaru {
	kame := &kamemaru{
		Echo: echo.New(),
	}
	err := envconfig.Process("turtle", &kame.Env)
	if err != nil {
		log.Fatal(err.Error())
	}
	kame.SetDeployMode()

	return kame
}

func (k *kamemaru) SetDeployMode() {
	switch k.Env.DeployMode {
	case "develop":
		k.SetLogger(os.Stderr)
	case "staging":
		logdir := os.Getenv("LOG_DIR")
		if logdir == "" {
			log.Fatal("LOG_DIR env was not set")
		}

		f, err := os.Create(filepath.Join(logdir, "laputa.log"))
		if err != nil {
			log.Fatal(err.Error())
		}
		defer f.Close()

		k.SetLogger(f)
	default:
		log.Fatal("LAPUTA_MODE env was not set")
	}

	k.Logger.Info("Graceful start laputa...", zap.Int("Port", l.env.Port))
}

func (k *kamemaru) SetLogger(Out zap.WriteSyncer) *zap.Logger {
	k.Logger = zap.New(
		zap.NewTextEncoder(zap.TextTimeFormat(time.ANSIC)),
		zap.AddCaller(), // Add line number option
		zap.Output(Out),
	)
}
