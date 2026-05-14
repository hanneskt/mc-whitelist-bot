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
}

func (b *Bot) Start(config *config.Config, logger *slog.Logger) (*bot.Client, error) {
	client, err := disgo.New(config.Token,
		bot.WithLogger(logger),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,       // which guild
				gateway.IntentGuildMembers, // join events
			),
		),

		// new member event
		bot.WithEventListenerFunc(b.OnMemberJoin),
		// slash command event
		bot.WithEventListenerFunc(HandleCommands),
		// modal submit event
		bot.WithEventListenerFunc(HandleModals),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize disgo client: %w", err)
	}

	// open discord
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {
		return nil, fmt.Errorf("Failed to connect to discord: %w", err)
	}

	// register commands
	_, err = client.Rest.SetGuildCommands(client.ApplicationID, config.GuildID, Commands())
	if err != nil {
		logger.Error("Error registering commands", "error", err)
	} else {
		logger.Info("Successfully registered commands!")
	}

	return client, nil
}

func (b *Bot) OnMemberJoin(e *events.GuildMemberJoin) {
	_, err := e.Client().Rest.CreateMessage(b.Config.WelcomeChannelID, discord.MessageCreate{
		Content: fmt.Sprintf("Welcome to the server, <@%s>!", e.Member.User.ID),
	})
	if err != nil {
		b.Logger.Error("Failed to send welcome message", "error", err)
	}
}
