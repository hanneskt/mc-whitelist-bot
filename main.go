package main

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"whitelistbot/config"
	"whitelistbot/discord"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		if errors.Is(err, config.ErrConfigCreated) {
			slog.Info("Generated empty config.json. Please fill it out and restart.")
			os.Exit(0)
		}

		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	discord.StartBot(cfg)

	// Keep the application running until you press CTRL+C
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
