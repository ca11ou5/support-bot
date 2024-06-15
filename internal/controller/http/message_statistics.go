package http

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
	"time"
)

type Stats struct {
	mu *sync.Mutex

	AllMessagesCount  int
	AllCommandsCount  int
	AllCallbacksCount int

	LatestMessagesCount  int
	LatestCommandsCount  int
	LatestCallbacksCount int
}

func (s *Server) HandleStats(update tgbotapi.Update) {
	ticker := time.NewTicker(time.Minute * 5)

	go func() {
		for range ticker.C {
			//INSERT STATS
			s.resetLatestStats()
		}
	}()
}

func (s *Server) statsCounting(update tgbotapi.Update) {
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

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
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

	s.stats.LatestMessagesCount = 0
	s.stats.LatestCommandsCount = 0
	s.stats.LatestCallbacksCount = 0
}

func (s *Server) saveStats() {
	s.useCase.
}
