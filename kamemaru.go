package kamemaru

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Code-Hex/kamemaru/internal/util"
	"github.com/Code-Hex/kamemaru/internal/validator"
	"github.com/Code-Hex/saltissimo"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/lestrrat/go-server-starter/listener"
	"github.com/uber-go/zap"
)

var (
	DeployMode string
	LogPath    string
)

// for config.toml
type (
	config struct {
		DB    db    `toml:"database"`
		Redis redis `toml:"redis"`
	}

	db struct {
		DBName   string `toml:"dbname"`
		Host     string `toml:"host"`
		UserName string `toml:"user"`
		Password string `toml:"pass"`
		Port     int    `toml:"port"`
		SSLmode  string `toml:"sslmode"`
	}

	redis struct {
		Host     string `toml:"host"`
		Network  string `toml:"network"`
		Password string `toml:"password"`
	}
)

// for kamemaru project
type kamemaru struct {
	Echo      *echo.Echo
	Logger    zap.Logger
	DB        *gorm.DB
	JWTSecret []byte
}

func New() *kamemaru {
	kame := &kamemaru{Echo: echo.New()}
	kame.Echo.Validator = validator.New()
	return kame.setup()
}

func (k *kamemaru) Run() int {
	defer k.DB.Close()

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
		err := os.MkdirAll(LogPath, os.ModeDir)
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

	var (
		err    error
		config config
	)
	// Skip this process if already exist "config.toml"
	if err = util.CreateConfig(); err != nil {
		log.Fatalf("Failed to create config toml: %s", err.Error())
	}

	if _, err = toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatalf("Failed to parse config toml: %s", err.Error())
	}

	if k.DB, err = gorm.Open("postgres", dbconf(config)); err != nil {
		log.Fatalf("Failed to connect database: %s", err.Error())
	}

	if err = k.route(config); err != nil {
		log.Fatalf("Failed to use middleware: %s", err.Error())
	}

	k.JWTSecret, err = saltissimo.RandomBytes(saltissimo.SaltLength)
	if err != nil {
		log.Fatalf("Failed to create secret: %s", err.Error())
	}

	return nil
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

func dbconf(conf config) string {
	var q string
	if conf.DB.Host != "" {
		q += "host=" + conf.DB.Host
	}

	if conf.DB.DBName != "" {
		q += " dbname=" + conf.DB.DBName
	}

	if conf.DB.UserName != "" {
		q += " user=" + conf.DB.UserName
	}

	if conf.DB.Password != "" {
		q += " password=" + conf.DB.Password
	}

	if conf.DB.SSLmode != "" {
		q += " sslmode=" + conf.DB.SSLmode
	}

	if conf.DB.Port > 0 {
		q += fmt.Sprintf(" port=%d", conf.DB.Port)
	}

	return q
}
