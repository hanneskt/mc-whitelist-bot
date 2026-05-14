package ptero

import (
	"context"
	"fmt"
	"log/slog"
	"whitelistbot/config"

	"github.com/davidarkless/go-pterodactyl"
	"github.com/davidarkless/go-pterodactyl/api"
)

type PteroClient struct {
	cfg    config.Config
	logger slog.Logger
	api    pterodactyl.Client
	ws     *api.WebsocketDetails
}

// make a new client for a pterodactyl panel
func New(cfg config.Config, logger slog.Logger) (*PteroClient, error) {
	client, err := pterodactyl.NewClient(cfg.PteroBaseURL, cfg.PteroApiKey, pterodactyl.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create pterodactyl client: %w", err)
	}

	ctx := context.Background()
	ws, err := client.ClientAPI.Servers(cfg.PteroServerIdentifier).GetWebsocket(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the pterodactyl websocket: %w", err)
	}

	return &PteroClient{
		cfg:    cfg,
		logger: logger,
		api:    *client,
		ws:     ws,
	}, nil
}

// whitelist a player on the server
func (c *PteroClient) WhitelistPlayer(username string) error {
	c.logger.Info("Whitelisting player", "username", username)
	command := fmt.Sprintf("whitelist add %s", username)
	err := c.api.ClientAPI.Servers(c.cfg.PteroServerIdentifier).SendCommand(context.Background(), command)
	if err != nil {
		return fmt.Errorf("failed to send whitelist command: %w", err)
	}

	c.logger.Info("Player whitelisted successfully", "username", username)
	return nil
}
