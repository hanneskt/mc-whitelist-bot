package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"whitelistbot/config"
	"whitelistbot/db"
	"whitelistbot/discord"
	"whitelistbot/ptero"
	"whitelistbot/service"

	_ "modernc.org/sqlite"
)

//go:embed sql/schema.sql
var ddl string

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

	baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	webhookHandler, err := discord.NewWebhookLogger(baseHandler, cfg.WebhookUrl)
	if err != nil {
		panic("Failed to create webhook logger: " + err.Error())
	}

	logger := slog.New(webhookHandler)

	pteroClient, err := ptero.New(*cfg, *logger)
	if err != nil {
		slog.Error("Failed to start Pterodactyl Client", "error", err)
	}

	ctx := context.Background()
	database, err := sql.Open("sqlite", "test.db?_pragma=foreign_keys(1)")
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	database.ExecContext(ctx, ddl)
	queries := db.New(database)

	minecraftSvc := service.NewMinecraftService(logger, pteroClient, queries)
	birthdaySvc := service.NewBirthdayService(logger, queries)

	// make bot
	bot := discord.Bot{
		Config:       cfg,
		Logger:       logger,
		MinecraftSvc: minecraftSvc,
		BirthdaySvc:  birthdaySvc,
	}

	err = bot.Start()
	if err != nil {
		logger.Error("Bot failed to start", "error", err)
		os.Exit(1)
	}
	defer bot.Stop()

	// Keep the application running until you press CTRL+C
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
