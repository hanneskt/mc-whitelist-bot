package discord

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"whitelistbot/service"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func WhitelistModal(serverName string) discord.ModalCreate {
	min_len := 3
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
					CustomID:    "minecraft_name",
					Style:       discord.TextInputStyleShort,
					MinLength:   &min_len, // this needs to be a pointer
					MaxLength:   16,
					Required:    true,
					Placeholder: "Santri",
				},
			},
			discord.LabelComponent{
				Label:       "Country",
				Description: "Where do you currently live?",
				Component: discord.TextInputComponent{
					CustomID:    "country",
					Style:       discord.TextInputStyleShort,
					Required:    true,
					MinLength:   &min_len,
					MaxLength:   32,
					Placeholder: "Germany",
				},
			},
			discord.LabelComponent{
				Label:       "Who invited you?",
				Description: "",
				Component: discord.TextInputComponent{
					CustomID:    "invited_by",
					Style:       discord.TextInputStyleShort,
					Required:    true,
					MinLength:   &min_len,
					MaxLength:   32,
					Placeholder: "Kameran",
				},
			},
		},
	}
}

var validNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (b *Bot) OnModalSubmit(e *events.ModalSubmitInteractionCreate) {
	b.Logger.Info("Handling a form submit", "form", e.Data.CustomID, "user", e.Member().EffectiveName())

	if e.Data.CustomID == "whitelist_form" {
		name := e.Data.Text("minecraft_name")
		if !validNameRegex.MatchString(name) {
			e.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf(`Your name "%s" doesn't seem te be a valid Minecraft name, can you check and try again?`, name),
			}.WithEphemeral(true))
			return
		}

		country := e.Data.Text("country")
		invited_by := e.Data.Text("invited_by")

		err := b.WhitelistSvc.WhitelistPlayer(service.PlayerToWhitelist{
			McUsername:  name,
			DiscordUuid: e.User().ID.String(),
			Country:     country,
			InvitedBy:   invited_by,
		})
		if err != nil {
			b.Logger.Error("Could not whitelist player", "error", err)

			// if the name was invalid
			if errors.Is(err, service.InvalidName) {
				e.CreateMessage(discord.MessageCreate{
					Content: fmt.Sprintf(`Your name "%s" doesn't seem te be a valid Minecraft name, can you check and try again?`, name),
				}.WithEphemeral(true))
				return
			}

			// otherwise
			e.CreateMessage(discord.MessageCreate{
				Content: "Something went wrong whitelisting you :(\nContact one of the admins (CEO role)",
			}.WithEphemeral(true))
			return
		}

		err = e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Hello %s, thanks for submitting the form!", name),
		})
		if err != nil {
			b.Logger.Error("Error replying to form submit", slog.Any("error", err))
		}
	}
}
