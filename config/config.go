package config

import (
	"encoding/json"
	"os"

	"github.com/disgoorg/snowflake/v2"
)

const configPath = "config.json"

type Config struct {
	Token            string       `json:"token"`
	GuildID          snowflake.ID `json:"guild_id"`
	WelcomeChannelID snowflake.ID `json:"welcome_channel_id"`
}

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
