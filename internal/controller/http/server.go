package http

import (
	"github.com/ca11ou5/support-bot/config"
	"github.com/ca11ou5/support-bot/internal/domain/message/entity"
	"github.com/ca11ou5/support-bot/internal/domain/message/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

type Server struct {
	useCase *usecase.UseCase

	stats entity.Stats
	bot   *tgbotapi.BotAPI
}

func NewServer(useCase *usecase.UseCase) *Server {
	return &Server{
		useCase: useCase,
		stats: entity.Stats{
			Timestamp: time.Now(),
			Words:     make(map[string]int),
		},
	}
}

func (s *Server) Start(cfg *config.Config) error {
	go func() error {
		_ = s.StartPolling(cfg.TelegramAPIToken)
		return nil
	}()

	err := s.RegisterHTTPServer()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
