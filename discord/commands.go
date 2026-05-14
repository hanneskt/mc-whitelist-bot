package discord

import (
	"log/slog"

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
			Name:        "hello",
			Description: "Says hello",
		},
	}
}

func (b *Bot) OnApplicationCommand(e *events.ApplicationCommandInteractionCreate) {
	b.Logger.Info("Handling the command", "command", e.Data.CommandName())
	switch e.Data.CommandName() {
	case "form":
		formCommand(e, b.Config.ServerName)
	case "hello":
		helloCommand(e)
	}
}

func formCommand(e *events.ApplicationCommandInteractionCreate, serverName string) {
	err := e.Modal(WhitelistModal(serverName))
	if err != nil {
		slog.Error("Error replying to form command", "error", err)
	}
}

func helloCommand(e *events.ApplicationCommandInteractionCreate) {

}
