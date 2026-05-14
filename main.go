package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"whitelistbot/config"
	"whitelistbot/discord"
)

func main() {
	logger := slog.Default()

	cfg, err := config.LoadConfig()
	if err != nil {
		if errors.Is(err, config.ErrConfigCreated) {
			slog.Info("Generated empty config.json. Please fill it out and restart.")
			os.Exit(0)
		}

		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	bot := discord.Bot{
		Config: cfg,
		Logger: logger,
	}
	client, err := bot.Start()
	if err != nil {
		logger.Error("Bot failed to start", "error", err)
		os.Exit(1)
	}
	defer client.Close(context.Background())

	// Keep the application running until you press CTRL+C
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
