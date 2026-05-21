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

type PlayerListResponse struct {
	Max    int `json:"max"`
	Online int `json:"online"`
	Sample []struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"sample"`
}

func (c *PteroManager) OnlinePlayers(ctx context.Context) ([]string, error) {
	server := c.servers[0] // just use the first server for this

	request, err := server.client.NewRequest(ctx, "GET", fmt.Sprintf("/api/client/servers/%s/minecraft-players", server.identifier), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for minecraft players: %w", err)
	}

	var players PlayerListResponse
	_, err = server.client.Do(ctx, request, &players)
	if err != nil {
		return nil, fmt.Errorf("failed to get minecraft players: %w", err)
	}

	names := make([]string, 0, len(players.Sample)) // an empty slice, but already allocated
	for _, p := range players.Sample {
		names = append(names, p.Name)
	}

	return names, nil
}
