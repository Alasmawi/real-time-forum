package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"reboot01.com/js/forum/internal/database"
	"reboot01.com/js/forum/internal/env"
	"reboot01.com/js/forum/internal/version"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

/*Define a config struct to hold all application configuration settings. 
values are set from env variables when the application starts.*/
type config struct {
	baseURL  string
	httpPort int
	db       struct {
		dsn         string
		// automigrate bool
	}
	jwt struct {
		secretKey string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	db     *database.DB
	logger *slog.Logger
	wg     sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 8080)
	cfg.db.dsn = env.GetString("DB_DSN", "db.sqlite")
	// cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.jwt.secretKey = env.GetString("JWT_SECRET_KEY", "rev3alim442itqpwlereeo5npf3h5uip")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.db.dsn/*, cfg.db.automigrate*/)
	if err != nil {
		return err
	}
	defer db.Close()

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
	}

	return app.serveHTTP()
}
