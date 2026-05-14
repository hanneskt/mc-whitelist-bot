package discord

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func WhitelistModal(serverName string) discord.ModalCreate {
	return discord.ModalCreate{
		CustomID: "whitelist_form",
		Title:    fmt.Sprintf("%s Whitelist Request", serverName),
		Components: []discord.LayoutComponent{
			discord.TextDisplayComponent{
				Content: "Fill in this form to get whitelisted on our **minecraft servers**!\n",
			},
			discord.LabelComponent{
				Label:       "Minecraft Username",
				Description: "The name you have in Minecraft",
				Component: discord.TextInputComponent{
					CustomID: "minecraft_name",
					Style:    discord.TextInputStyleShort,
				},
			},
			discord.LabelComponent{
				Label:       "Country",
				Description: "Where do you currently live?",
				Component: discord.TextInputComponent{
					CustomID: "country",
					Style:    discord.TextInputStyleShort,
				},
			},
			discord.LabelComponent{
				Label:       "Who invited you?",
				Description: "",
				Component: discord.TextInputComponent{
					CustomID: "who_invite",
					Style:    discord.TextInputStyleShort,
				},
			},
		},
	}
}

func (b *Bot) OnModalSubmit(e *events.ModalSubmitInteractionCreate) {
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
