package usecase

import (
	"github.com/ca11ou5/support-bot/internal/domain/message/repository"
)

type UseCase struct {
	messageRepo *repository.MessageRepository
}

func NewMessageUseCase(messageRepo *repository.MessageRepository) *UseCase {
	return &UseCase{messageRepo: messageRepo}
}

func (uc *UseCase) HandleCommand(command string) string {

	// TODO: logic if commandAction.needsStateChange true
	action := uc.messageRepo.GetCommandAction(command)

	return action.Text
}
