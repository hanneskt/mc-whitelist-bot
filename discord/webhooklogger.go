package discord

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
)

type WebhookLogger struct {
	slog.Handler
	webhookClient *webhook.Client
}

func NewWebhookLogger(baseHandler slog.Handler, webhookURL string) (*WebhookLogger, error) {
	client, err := webhook.NewWithURL(webhookURL)
	if err != nil {
		return nil, err
	}

	return &WebhookLogger{
		Handler:       baseHandler,
		webhookClient: client,
	}, nil
}

func (w *WebhookLogger) Handle(ctx context.Context, r slog.Record) error {
	err := w.Handler.Handle(ctx, r)

	if r.Level >= slog.LevelInfo {
		go func(record slog.Record) {
			content := fmt.Sprintf("**[%s]** %s", record.Level.String(), record.Message)

			var attrs string
			record.Attrs(func(a slog.Attr) bool {
				attrs += fmt.Sprintf("%s: %v\n", a.Key, a.Value.Any())
				return true
			})

			if attrs != "" {
				content += fmt.Sprintf("\n```yaml\n%s\n```", attrs)
			}

			_, _ = w.webhookClient.CreateMessage(discord.WebhookMessageCreate{
				Content: content,
			}, rest.CreateWebhookMessageParams{})
		}(r)
	}

	return err
}
