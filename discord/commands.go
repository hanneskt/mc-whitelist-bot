package discord

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "form",
			Description: "Opens the application form",
		},
		discord.SlashCommandCreate{
			Name:        "players",
			Description: "Lists players on the Survival server",
		},
		discord.SlashCommandCreate{
			Name:        "hello",
			Description: "Says hello",
		},
	}
}

func (b *Bot) OnApplicationCommand(e *events.ApplicationCommandInteractionCreate) {
	b.Logger.Info("Handling the command", "command", e.Data.CommandName(), "user", e.Member().EffectiveName())
	var err error

	switch e.Data.CommandName() {
	case "form":
		err = b.formCommand(e)
	case "players":
		err = b.playersCommand(e)
	case "hello":
		err = b.helloCommand(e)
	}
	if err != nil {
		slog.Error("Error replying to command", "command", e.Data.CommandName(), "error", err)
	}
}

func (b *Bot) formCommand(e *events.ApplicationCommandInteractionCreate) error {
	return e.Modal(WhitelistModal(b.Config.DiscordServerName))
}

func (b *Bot) playersCommand(e *events.ApplicationCommandInteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	players, err := b.MinecraftSvc.OnlinePlayers(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder
	for _, player := range players {
		fmt.Fprintf(&sb, "• %s\n", player)
	}

	e.CreateMessage(discord.NewMessageCreate().AddEmbeds(discord.Embed{
		Title:       fmt.Sprintf("Survival Online Players (%d)", len(players)),
		Description: sb.String(),
		Color:       0x00FF00,
	}).WithEphemeral(true))

	return nil
}

func (b *Bot) helloCommand(e *events.ApplicationCommandInteractionCreate) error {
	e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Hello <@%s>", e.User().ID),
	})

	return nil
}
