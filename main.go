package main

import (
	"awesomeProject/Internal/config"
	"fmt"
	"log"
	"log/slog"
	"os"

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
	fmt.Println(cfg)
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
