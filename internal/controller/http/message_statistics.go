package http

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log/slog"
	"os"
	"time"
)

func (s *Server) HandleStats() {
	ticker := time.NewTicker(time.Minute * 1)

	go func() {
		for range ticker.C {
			err := s.useCase.SaveStats(s.stats)
			if err != nil {
				slog.Error(err.Error())
			}

			s.resetLatestStats()
		}
	}()
}

func (s *Server) statsCounting(update tgbotapi.Update) {
	s.stats.Mu.Lock()
	defer s.stats.Mu.Unlock()

	s.stats.AllMessagesCount++
	s.stats.LatestMessagesCount++

	if update.Message.IsCommand() {
		s.stats.AllCommandsCount++
		s.stats.LatestCommandsCount++
	} else if update.CallbackQuery != nil {
		s.stats.AllCallbacksCount++
		s.stats.LatestCallbacksCount++
	}
}

func (s *Server) resetLatestStats() {
	s.stats.Mu.Lock()
	defer s.stats.Mu.Unlock()

	s.stats.LatestMessagesCount = 0
	s.stats.LatestCommandsCount = 0
	s.stats.LatestCallbacksCount = 0

	s.stats.Timestamp = time.Now()
}

func (s *Server) saveStats() {
	err := s.useCase.SaveStats(s.stats)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
