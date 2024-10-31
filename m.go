package main

import (
	"GoCRUD/api"
	"GoCRUD/database"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		os.Exit(1)
	}

	slog.Info("all systems offline")
}

func run() error {
	handle := api.NewHandler()
	database.InitializeDatabase()

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handle,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
