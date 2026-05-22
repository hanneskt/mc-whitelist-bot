package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
)

func (b *Bot) StartBirthdayAnnouncer() {
	go func() {
		b.Logger.Info("Birthday scheduler loop started")
		for {
			now := time.Now()

			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
			if now.After(nextRun) {
				nextRun = nextRun.AddDate(0, 0, 1)
			}

			delay := time.Until(nextRun)
			b.Logger.Info("Next birthday check scheduled", "at", nextRun.Format("2006-01-02 15:04:05"), "delay", delay.String())

			timer := time.NewTimer(delay)
			<-timer.C

			b.Logger.Info("Running daily birthday check")
			b.announceBirthdays()
		}
	}()
}

func (b *Bot) announceBirthdays() {
	now := time.Now()
	birthdays, err := b.BirthdaySvc.GetBirthdaysForDate(now.Day(), int(now.Month()))
	if err != nil {
		b.Logger.Error("Failed to fetch today's birthdays", "error", err)
		return
	}

	if len(birthdays) == 0 {
		b.Logger.Info("No one has a birthday today")
		return
	}

	var mentions []string
	for _, bday := range birthdays {
		mentions = append(mentions, fmt.Sprintf("<@%s>", bday.ID))
	}

	message := fmt.Sprintf("🎉 Happy Birthday to %s! 🎂", strings.Join(mentions, " and "))

	b.botClient.Rest.CreateMessage(b.Config.BirthdayChannelID, discord.NewMessageCreate().WithContent(message))
}
