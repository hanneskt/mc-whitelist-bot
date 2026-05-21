package discord

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"whitelistbot/service"

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
		discord.SlashCommandCreate{
			Name:        "birthday",
			Description: "Commands relating to the birthdaycommands",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "set",
					Description: "Set you birthday so the bot can announce your birthday!",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "day",
							Description: "The day of the month of your birthday",
							Required:    true,
							MinValue:    intPtr(1),
							MaxValue:    intPtr(31),
						},
						discord.ApplicationCommandOptionInt{
							Name:        "month",
							Description: "The month of your birthday, as a number",
							Required:    true,
							MinValue:    intPtr(1),
							MaxValue:    intPtr(12),
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "get",
					Description: "Get the birthdays for a month",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "month",
							Description: "The month to get birthdays for",
							MinValue:    intPtr(1),
							MaxValue:    intPtr(12),
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "remove",
					Description: "Remove your birthday",
				},
			},
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
	case "birthday":
		err = b.birthdayCommand(e)
	default:
		err = errors.New("No handler for this command")
	}
	if err != nil {
		e.CreateMessage(discord.NewMessageCreate().WithContent("There was an error running the command :(").WithEphemeral(true))
		b.Logger.Error("Error replying to command", "command", e.Data.CommandName(), "error", err)
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

func (b *Bot) birthdayCommand(e *events.ApplicationCommandInteractionCreate) error {
	if e.SlashCommandInteractionData().SubCommandName == nil {
		return errors.New("birthday subcommand required")
	}

	switch *e.SlashCommandInteractionData().SubCommandName {
	case "set":
		day := e.SlashCommandInteractionData().Int("day")
		month := e.SlashCommandInteractionData().Int("month")

		err := b.BirthdaySvc.SetBirthday(e.User().ID.String(), day, month)
		if err != nil {
			if errors.Is(err, service.ErrBirthdayAlreadySet) {
				return e.CreateMessage(discord.NewMessageCreate().
					WithContent("You already set your birthday! You can remove it with /birthday remove").
					WithEphemeral(true))
			}
			return err
		}

		e.CreateMessage(discord.NewMessageCreate().WithContent("Your birthday was saved!").WithEphemeral(true))
	case "get":
		month := e.SlashCommandInteractionData().Int("month")

		if month == 0 {
			month = int(time.Now().Month())
		}

		birthdays, err := b.BirthdaySvc.GetBirthdaysForMonth(month)
		if err != nil {
			return err
		}

		var sb strings.Builder
		if len(birthdays) == 0 {
			sb.WriteString("_No birthdays this month._")
		} else {
			for _, bday := range birthdays {
				fmt.Fprintf(&sb, "**%d** - <@%s>\n", bday.Day, bday.ID)
			}
		}

		monthName := time.Month(month).String()
		return e.CreateMessage(discord.NewMessageCreate().AddEmbeds(discord.Embed{
			Title:       fmt.Sprintf("Birthdays in %s", monthName),
			Description: sb.String(),
			Color:       0x00FF00,
		}).WithEphemeral(true))
	case "remove":
		if err := b.BirthdaySvc.DeleteBirthday(e.User().ID.String()); err != nil {
			return err
		}
		e.CreateMessage(discord.NewMessageCreate().WithContent("Your birthday was deleted!").WithEphemeral(true))
	default:
		return errors.New("no such birthday subcommand")
	}

	return nil
}

func intPtr(i int) *int {
	return &i
}
