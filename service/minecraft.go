package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"whitelistbot/db"
	"whitelistbot/ptero"
)

type MinecraftService struct {
	logger  *slog.Logger
	ptero   *ptero.PteroManager
	queries *db.Queries
}

var InvalidName = errors.New("invalid minecraft username")
var AlreadyWhitelisted = errors.New("already whitelisted")

func NewMinecraftService(l *slog.Logger, p *ptero.PteroManager, q *db.Queries) *MinecraftService {
	return &MinecraftService{
		logger:  l,
		ptero:   p,
		queries: q,
	}
}

type PlayerToWhitelist struct {
	McUsername  string
	DiscordUuid string
	Country     string
	InvitedBy   string
}

func (s *MinecraftService) WhitelistPlayer(player PlayerToWhitelist) error { // TODO: return a whitelist result
	_, err := s.queries.GetPlayerByDiscordUuid(context.Background(), player.DiscordUuid)
	if err == nil {
		return AlreadyWhitelisted
	}

	playerInfo, err := s.UsernameValid(player.McUsername)
	if err != nil {
		return err
	}

	err = s.ptero.WhitelistPlayerCommand(playerInfo.Name)
	if err != nil {
		return fmt.Errorf("sending whitelist command to server failed: %w", err)
	}

	dbPlayer, err := s.queries.CreatePlayer(context.Background(), db.CreatePlayerParams{
		McUuid:      playerInfo.Uuid,
		McUsername:  playerInfo.Name,
		DiscordUuid: player.DiscordUuid,
		Country:     player.Country,
		InvitedBy:   player.InvitedBy,
		Whitelisted: true,
	})
	if err != nil {
		return fmt.Errorf("failed to save player to database: %w", err)
	}
	s.logger.Info("Saved player to the database", "player", dbPlayer)

	return nil
}

type MojangPlayerInfo struct {
	Uuid string `json:"id"`
	Name string `json:"name"`
}

func (s *MinecraftService) UsernameValid(username string) (*MojangPlayerInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", username))
	if err != nil {
		s.logger.Warn("Mojang api returned an error", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body of mojang request")
		}

		var playerInfo MojangPlayerInfo
		json.Unmarshal(body, &playerInfo)
		return &playerInfo, nil
	}

	if resp.StatusCode == 404 {
		s.logger.Warn("Username not found", "username", username)
		return nil, InvalidName
	}

	return nil, fmt.Errorf("mojang api returned status code: %d", resp.StatusCode)
}

func (s *MinecraftService) OnlinePlayers(ctx context.Context) ([]string, error) {
	return s.ptero.OnlinePlayers(ctx)
}
