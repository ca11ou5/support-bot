package usecase

import (
	"github.com/ca11ou5/support-bot/internal/domain/message/entity"
)

func (uc *UseCase) SaveStats(stats entity.Stats) error {
	return uc.messageRepo.SaveStats(stats)
}

func (uc *UseCase) InsertWords(words map[string]int) error {
	return uc.messageRepo.InsertWords(words)
}
