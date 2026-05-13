package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"whitelistbot/config"
	"whitelistbot/discord"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	discord.StartBot(config)

	// Keep the application running until you press CTRL+C
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
