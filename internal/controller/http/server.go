package http

import (
	"github.com/ca11ou5/support-bot/config"
	"github.com/ca11ou5/support-bot/internal/domain/message/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Server struct {
	useCase *usecase.UseCase

	stats *Stats
	bot   *tgbotapi.BotAPI
}

func NewServer(useCase *usecase.UseCase) *Server {
	return &Server{
		useCase: useCase,
	}
}

func (s *Server) Start(cfg *config.Config) error {
	err := s.StartPolling(cfg.TelegramAPIToken)
	if err != nil {
		return err
	}

	return nil
}
