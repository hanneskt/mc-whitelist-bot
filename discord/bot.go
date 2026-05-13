package discord

import (
	"context"
	"fmt"
	"log/slog"
	"whitelistbot/config"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func StartBot(config *config.Config) {
	client, err := disgo.New(config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,       // which guild
				gateway.IntentGuildMembers, // join events
			),
		),

		// new member event
		bot.WithEventListenerFunc(func(e *events.GuildMemberJoin) {
			_, err := e.Client().Rest.CreateMessage(config.WelcomeChannelID, discord.MessageCreate{
				Content: fmt.Sprintf("Welcome to the server, <@%s>!", e.Member.User.ID),
			})
			if err != nil {
				slog.Error("Failed to send welcome message", slog.Any("err", err))
			}
		}),

		// slash command event
		bot.WithEventListenerFunc(HandleCommands),
		bot.WithEventListenerFunc(HandleModals),
	)
	if err != nil {
		panic(err)
	}

	// open discord
	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	// register commands
	_, err = client.Rest.SetGuildCommands(client.ApplicationID, 193467328661422081, Commands())
	if err != nil {
		slog.Error("Error registering commands", slog.Any("err", err))
	} else {
		slog.Info("Successfully registered commands!")
	}

}
