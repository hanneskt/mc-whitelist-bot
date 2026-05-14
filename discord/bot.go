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

func (self *Bot) Start() (*bot.Client, error) {
	client, err := disgo.New(self.Config.Token,
		bot.WithLogger(self.Logger),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,       // which guild
				gateway.IntentGuildMembers, // join events
			),
		),
		bot.WithEventListenerFunc(self.OnMemberJoin),
		bot.WithEventListenerFunc(self.OnApplicationCommand),
		bot.WithEventListenerFunc(self.OnModalSubmit),
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
	_, err = client.Rest.SetGuildCommands(client.ApplicationID, self.Config.GuildID, Commands())
	if err != nil {
		self.Logger.Error("Error registering commands", "error", err)
	} else {
		self.Logger.Info("Successfully registered commands!")
	}

	return client, nil
}

func (self *Bot) OnMemberJoin(e *events.GuildMemberJoin) {
	_, err := e.Client().Rest.CreateMessage(self.Config.WelcomeChannelID, discord.MessageCreate{
		Content: fmt.Sprintf("Welcome to the server, <@%s>!", e.Member.User.ID),
	})
	if err != nil {
		self.Logger.Error("Failed to send welcome message", "error", err)
	}
}
