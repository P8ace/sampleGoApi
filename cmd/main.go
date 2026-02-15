package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {

	//setup structured logging
	logHandlerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, logHandlerOptions))
	slog.SetDefault(log)

	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
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
