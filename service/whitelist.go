package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"whitelistbot/ptero"
)

type WhitelistService struct {
	logger *slog.Logger
	ptero  *ptero.PteroClient
}

var InvalidName = errors.New("invalid minecraft username")

func NewWhitelistService(l *slog.Logger, p *ptero.PteroClient) *WhitelistService {
	return &WhitelistService{
		logger: l,
		ptero:  p,
	}
}

func (s *WhitelistService) WhitelistPlayer(username string) error { // TODO: return a whitelist result
	if !s.UsernameValid(username) {
		return InvalidName
	}

	err := s.ptero.WhitelistPlayerCommand(username)
	if err != nil {
		return fmt.Errorf("sending whitelist command to server failed: %w", err)
	}

	return nil
}

func (s *WhitelistService) UsernameValid(username string) bool {
	resp, err := http.Get(fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", username))
	if err != nil {
		s.logger.Warn("Mojang api returned an error", "error", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		s.logger.Warn("Username not found", "username", username)
		return false
	}

	if resp.StatusCode == 200 {
		return true
	}

	return false
}
