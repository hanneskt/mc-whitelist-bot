package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"whitelistbot/db"
)

type BirthdayService struct {
	logger  *slog.Logger
	queries *db.Queries
}

func NewBirthdayService(l *slog.Logger, q *db.Queries) *BirthdayService {
	return &BirthdayService{
		logger:  l,
		queries: q,
	}
}

var ErrBirthdayAlreadySet = errors.New("birthday already set")

func (s *BirthdayService) SetBirthday(discordId string, day int, month int) error {
	birthday, err := s.GetBirthdayById(discordId)
	if birthday != nil {
		return ErrBirthdayAlreadySet
	}

	_, err = s.queries.CreateBirthday(context.Background(), db.CreateBirthdayParams{
		DiscordUuid: discordId,
		Day:         int64(day),
		Month:       int64(month),
	})
	if err != nil {
		return fmt.Errorf("failed to save birthday: %w", err)
	}

	return nil
}

type Birthday struct {
	ID    string
	Day   int
	Month int
}

func (s *BirthdayService) GetBirthdayById(discordId string) (*Birthday, error) {
	birthday, err := s.queries.GetBirthdayById(context.Background(), discordId)
	if err != nil {
		return nil, fmt.Errorf("failed to get birthday: %w", err)
	}

	return &Birthday{
		ID:    birthday.DiscordUuid,
		Day:   int(birthday.Day),
		Month: int(birthday.Month),
	}, nil
}

func (s *BirthdayService) DeleteBirthday(discordId string) error {
	err := s.queries.DeleteBirthday(context.Background(), discordId)
	if err != nil {
		return fmt.Errorf("failed to delete birthday: %w", err)
	}

	return nil
}

func (s *BirthdayService) GetBirthdaysForDate(day, month int) error {
	_, err := s.queries.GetBirthdaysForDate(context.Background(), db.GetBirthdaysForDateParams{
		Day:   int64(day),
		Month: int64(month),
	})
	if err != nil {
		return fmt.Errorf("failed to get birthdays for date %d/%d: %w", day, month, err)
	}

	return nil
}

func (s *BirthdayService) GetBirthdaysForMonth(month int) ([]Birthday, error) {
	dbBirthdays, err := s.queries.GetBirthdaysForMonth(context.Background(), int64(month))
	if err != nil {
		return nil, fmt.Errorf("failed to get birthdays for month %d: %w", month, err)
	}

	birthdays := make([]Birthday, 0, len(dbBirthdays)) // an empty slice, but already allocated
	for _, b := range dbBirthdays {
		birthdays = append(birthdays, Birthday{
			ID:    b.DiscordUuid,
			Day:   int(b.Day),
			Month: int(b.Month),
		})
	}

	return birthdays, nil
}
