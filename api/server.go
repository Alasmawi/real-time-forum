package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"reboot01.com/js/realtime-forum/internal/database"
)

/*
Define a config struct to hold all Application configuration settings.
values are set from env variables when the Application starts.
*/
type Config struct {
	BaseURL  string
	HttpPort int
	DB       struct {
		DSN string
		// automigrate bool
	}
	JWT struct {
		SecretKey string
	}
}

// Define an Application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type Application struct {
	Config Config
	DB     *database.DB
	Logger *slog.Logger
	WG     sync.WaitGroup
}

// neuteredFileSystem is a custom type which embeds the standard http.FileSystem.
type neuteredFileSystem struct {
	fs http.FileSystem
}

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (app *Application) ServeHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.HttpPort),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	app.Logger.Info("starting server", slog.Group("server", "addr", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	app.Logger.Info("stopped server", slog.Group("server", "addr", srv.Addr))

	app.WG.Wait()
	return nil
}

// Open checks the requested file path and determines whether it is a directory or not.
// If it is a directory, it attempts to open an index.html file in it.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			// Closes the original file to avoid a file descriptor leak
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
