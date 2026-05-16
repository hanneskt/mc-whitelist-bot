package ptero

import (
	"context"
	"fmt"
	"log/slog"
	"whitelistbot/config"

	"github.com/davidarkless/go-pterodactyl"
)

type PteroClient struct {
	cfg    config.Config
	logger slog.Logger
	api    pterodactyl.Client
}

// make a new client for a pterodactyl panel
func New(cfg config.Config, logger slog.Logger) (*PteroClient, error) {
	client, err := pterodactyl.NewClient(cfg.PteroBaseURL, cfg.PteroApiKey, pterodactyl.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create pterodactyl client: %w", err)
	}

	return &PteroClient{
		cfg:    cfg,
		logger: logger,
		api:    *client,
	}, nil
}

// whitelist a player on the server
func (c *PteroClient) WhitelistPlayerCommand(username string) error {
	c.logger.Info("Whitelisting player", "username", username)
	command := fmt.Sprintf("whitelist add %s", username)
	err := c.api.ClientAPI.Servers(c.cfg.PteroServerIdentifier).SendCommand(context.Background(), command)
	if err != nil {
		return fmt.Errorf("failed to send whitelist command: %w", err)
	}

	c.logger.Info("Player whitelisted successfully", "username", username)
	return nil
}
