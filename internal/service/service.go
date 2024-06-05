package service

import "github.com/ca11ou5/support-bot/internal/repository"

type Service struct {
	AdminRepository
	TelegramRepository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		AdminRepository: *repo,
	}
}

type AdminRepository interface{}

type TelegramRepository interface {
}
