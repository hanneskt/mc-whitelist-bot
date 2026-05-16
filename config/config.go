package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/disgoorg/snowflake/v2"
)

var ErrConfigCreated = errors.New("config file created; please edit and restart")

const configPath = "config.json"

type Config struct {
	WebhookUrl        string       `json:"webhook_url"`
	Token             string       `json:"token"`
	GuildID           snowflake.ID `json:"guild_id"`
	WelcomeChannelID  snowflake.ID `json:"welcome_channel_id"`
	DiscordServerName string       `json:"discord_server_name"`

	PteroServers []PteroServer `json:"servers"`
}

type PteroServer struct {
	PteroBaseURL string `json:"pterodactyl_base_url"`
	PteroApiKey  string `json:"pterodactyl_api_key"`

	Name                  string `json:"name"`
	PteroServerIdentifier string `json:"pterodactyl_server_identifier"`
}

func LoadConfig() (*Config, error) {
	var config Config

	file, err := os.ReadFile(configPath)

	// if there is no file, generate an empty config
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			emptyConfig, _ := json.MarshalIndent(config, "", "  ")

			err = os.WriteFile(configPath, emptyConfig, 0600)
			if err != nil {
				return nil, fmt.Errorf("Failed to write config: %w", err)
			}
			return nil, ErrConfigCreated
		}
		return nil, fmt.Errorf("Failed to read config: %w", err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("Failed to read config json: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	var errs []error

	if len(c.Token) != 72 {
		errs = append(errs, errors.New("token is missing"))
	}

	if c.GuildID == 0 {
		errs = append(errs, errors.New("guild_id cannot be 0"))
	}

	if c.WelcomeChannelID == 0 {
		errs = append(errs, errors.New("welcome_channel_id cannot be 0"))
	}

	for _, server := range c.PteroServers {
		if server.PteroBaseURL == "" {
			errs = append(errs, errors.New("pterodactyl_base_url shouldn't be empty"))
		}
		if server.PteroApiKey == "" {
			errs = append(errs, errors.New("pterodactyl_api_key shouldn't be empty"))
		}
		if server.PteroServerIdentifier == "" {
			errs = append(errs, errors.New("pterodactyl_server_identifier shouldn't be empty"))
		}
	}

	return errors.Join(errs...)
}
