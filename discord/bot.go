package discord

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"whitelistbot/config"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type Bot struct {
	Config *config.Config
	Logger *slog.Logger

	client *bot.Client
}

func (b *Bot) Start() error {
	var err error

	b.client, err = disgo.New(b.Config.Token,
		bot.WithLogger(b.Logger),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,       // which guild
				gateway.IntentGuildMembers, // join events
			),
		),
		bot.WithEventListenerFunc(b.OnMemberJoin),
		bot.WithEventListenerFunc(b.OnApplicationCommand),
		bot.WithEventListenerFunc(b.OnModalSubmit),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize disgo client: %w", err)
	}

	// open discord
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = b.client.OpenGateway(ctx); err != nil {
		return fmt.Errorf("failed to connect to discord: %w", err)
	}

	// register commands
	_, err = b.client.Rest.SetGuildCommands(b.client.ApplicationID, b.Config.GuildID, Commands())
	if err != nil {
		b.Logger.Error("Error registering commands", "error", err)
	} else {
		b.Logger.Info("Successfully registered commands!")
	}

	return nil
}

func (b *Bot) Stop(ctx context.Context) {
	if b.client != nil {
		b.Logger.Info("Shutting down bot client")
		b.client.Close(ctx)
	}
}

func (b *Bot) OnMemberJoin(e *events.GuildMemberJoin) {
	_, err := e.Client().Rest.CreateMessage(b.Config.WelcomeChannelID, discord.MessageCreate{
		Content: fmt.Sprintf("Welcome to the server, <@%s>!", e.Member.User.ID),
	})
	if err != nil {
		b.Logger.Error("Failed to send welcome message", "error", err)
	}
}
