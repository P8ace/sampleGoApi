package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/P8ace/sampleGoApi/internal/adapters/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"),
		},
	}

	//setup structured logging
	logHandlerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, logHandlerOptions))
	slog.SetDefault(log)

	// Database connection
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("Error while connecting to database", "Error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	slog.Info("Connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	handler := api.mount()
	server := api.getServer((handler))

	runGroup := Group{}
	runGroup.Add(func() error {
		slog.Info("Starting HTTP Server", "Address", api.config.addr)
		return server.ListenAndServe()
	}, func(err error) {
		server.Shutdown(context.Background())
	})
	runGroup.Add(listenInterrupts(context.Background()))

	initialError := runGroup.Run()
	slog.Error("Service terminated", "Root Cause", initialError)
}
