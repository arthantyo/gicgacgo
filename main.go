package main

import (
	shared "gicgacgo/shared"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("error loading .env file")
	}

	// Load stats from JSON file
	if err := shared.LoadStats(); err != nil {
		slog.Error("error loading stats", slog.Any("error", err))
	} else {
		slog.Info("stats loaded successfully")
	}

	dg, err := setupBot()
	if err != nil {
		slog.Error("error setting up bot,")
		slog.Any("error", err)
	}

	dg.AddHandler(handleInteraction)
	dg.AddHandler(handleReady)

	if err := dg.Open(); err != nil {
		slog.Error("cannot open the session: ", slog.Any("error", err))

	}

	slog.Info("adding commands...")
	registeredCommands, err := registerCommands(dg)
	if err != nil {
		slog.Error("error registering commands", slog.Any("error", err))

	}

	defer dg.Close()

	slog.Info("bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	slog.Info("removing commands...")
	removeCommands(dg, registeredCommands)

	slog.Info("gracefully shutting down.")
}
