package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"reboot01.com/js/realtime-forum/api"
	"reboot01.com/js/realtime-forum/internal/database"
	"reboot01.com/js/realtime-forum/internal/env"
	"reboot01.com/js/realtime-forum/internal/version"
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

func run(logger *slog.Logger) error {
	var cfg api.Config

	cfg.BaseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.HttpPort = env.GetInt("HTTP_PORT", 8080)
	cfg.DB.DSN = env.GetString("DB_DSN", "db.sqlite")
	// cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	// cfg.JWT.SecretKey = env.GetString("JWT_SECRET_KEY", "rev3alim442itqpwlereeo5npf3h5uip")

	showVersion := flag.Bool("version", false, "display version and exit")
	// seedDB := flag.Bool("seed", true, "seed the database with sample data")
	// clearDB := flag.Bool("clear", false, "clear all data from the database")
	// reseedDB := flag.Bool("reseed", false, "clear and reseed the database")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.DB.DSN /*, cfg.db.automigrate*/)
	if err != nil {
		return err
	}
	defer db.Close()

	// if *clearDB {
	// 	return db.ClearDatabase()
	// }

	// if *reseedDB {
	// 	if err := db.ClearDatabase(); err != nil {
	// 		return err
	// 	}
	// 	return db.SeedData()
	// }

	// if *seedDB {
	// 	return db.SeedIfEmpty()
	// }

	app := &api.Application{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}

	return app.ServeHTTP()
}
