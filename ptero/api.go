package ptero

import (
	"context"
	"fmt"
	"log/slog"
	"whitelistbot/config"

	"github.com/davidarkless/go-pterodactyl"
)

type PteroManager struct {
	cfg    config.Config
	logger slog.Logger

	servers []Server
}

type Server struct {
	name       string
	identifier string
	client     *pterodactyl.Client
}

// make a new client for a pterodactyl panel
func New(cfg config.Config, logger slog.Logger) (*PteroManager, error) {
	var pteroClients []Server
	for _, server := range cfg.PteroServers {
		client, err := pterodactyl.NewClient(server.PteroBaseURL, server.PteroApiKey, pterodactyl.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create pterodactyl client: %w", err)
		}

		pteroClients = append(pteroClients, Server{
			name:       server.Name,
			identifier: server.PteroServerIdentifier,
			client:     client,
		})
	}

	return &PteroManager{
		cfg:     cfg,
		logger:  logger,
		servers: pteroClients,
	}, nil
}

// whitelist a player on the server
func (c *PteroManager) WhitelistPlayerCommand(username string) error {
	c.logger.Info("Whitelisting player", "username", username)
	command := fmt.Sprintf("whitelist add %s", username)

	for _, server := range c.servers {
		c.logger.Info("Whitelisting player on server", "server", server.name)
		err := server.client.ClientAPI.Servers(server.identifier).SendCommand(context.Background(), command)
		if err != nil {
			return fmt.Errorf("failed to send whitelist command: %w", err)
		}
	}

	c.logger.Info("Player whitelisted successfully", "username", username)
	return nil
}
