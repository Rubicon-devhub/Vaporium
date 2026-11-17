package main

import (
	"awesomeProject/Internal/config"
	midlLog "awesomeProject/Internal/http-server/midlware/loger"
	"awesomeProject/Internal/lib/handlers/url"
	"awesomeProject/Internal/lib/sl"
	"net/http"

	postgresql "awesomeProject/storage/psql"
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
const (
	envlocal = "local"
	envdev   = "dev"
	envprod  = "prod"
)

func main() {
	err := godotenv.Load()

	cfg := config.MustLoad()

	if err != nil {
		log.Fatalf("cannot load env: %s", err)
	}

	loger := setupLoger(cfg.Env)

	loger.Debug("Loger initialized", slog.String("env", cfg.Env))

	ctx, cansel := context.WithTimeout(context.Background(), cfg.DbTimeout)

	defer cansel()

	storage, err := postgresql.New(ctx, cfg.ConnectionString)
	if err != nil {
		loger.Error("failed to init storage", sl.Err(err))
		return
	}
	loger.Debug("Data base initializated", slog.String("env", cfg.Env))

	router := chi.NewRouter()

	mv := midlLog.New(loger)

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	router.Use(mv)

	router.Post("/url", url.New(loger, storage, ctx))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err:= srv.ListenAndServe(); err!=nil{
		loger.Error("failed to start server")
		
	}
}

func setupLoger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envlocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envprod:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
