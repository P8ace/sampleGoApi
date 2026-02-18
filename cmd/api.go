package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/P8ace/sampleGoApi/internal/adapters/postgres/repo"
	"github.com/P8ace/sampleGoApi/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID) // rate limiting
	r.Use(middleware.RealIP)    //import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //recover from crashes
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All good"))
	})

	productService := products.NewService(repo.New(app.db))
	productsHandler := products.NewHandler(productService)

	r.Get("/products", productsHandler.ListProducts)

	return r
}

func (app *application) getServer(h http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv
}

type dbConfig struct {
	dsn string
}
type config struct {
	addr string
	db   dbConfig
}

// function to handle graceful shutdowns
type SignalError struct {
	sig os.Signal
}

// Error implements the error interface.
func (e SignalError) Error() string {
	return fmt.Sprintf("received signal %s", e.sig)
}

func listenInterrupts(ctx context.Context) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	return func() error {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer signal.Stop(sigChan)
			select {
			case sig := <-sigChan:
				return SignalError{sig: sig}
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(err error) {
			cancel()
		}
}
