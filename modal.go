package main

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func WhitelistModal() discord.ModalCreate {
	return discord.ModalCreate{
		CustomID: "whitelist_form",
		Title:    "Whitelist Request",
		Components: []discord.LayoutComponent{
			discord.LabelComponent{
				Label:       "Minecraft Username",
				Description: "The name you have in Minecraft",
				Component: discord.TextInputComponent{
					CustomID: "minecraft_name",
					Style:    discord.TextInputStyleShort,
				},
			},
		},
	}
}

func HandleModals(e *events.ModalSubmitInteractionCreate) {
	if e.Data.CustomID == "whitelist_form" {
		name := e.Data.Text("minecraft_name")

		err := e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Hello %s, thanks for submitting the form!", name),
		})
		if err != nil {
			slog.Error("Error replying to form submit", slog.Any("error", err))
		}
	}
}
