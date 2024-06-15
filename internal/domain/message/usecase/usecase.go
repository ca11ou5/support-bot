package usecase

import (
	"github.com/ca11ou5/support-bot/internal/domain/message/repository"
	"strconv"
)

type UseCase struct {
	messageRepo *repository.MessageRepository
}

func NewMessageUseCase(messageRepo *repository.MessageRepository) *UseCase {
	return &UseCase{messageRepo: messageRepo}
}

func (uc *UseCase) HandleCommand(command string, chatID int64) string {
	action := uc.messageRepo.GetCommandAction(command)
	id := strconv.Itoa(int(chatID))

	if action.NeedsSetupState {
		uc.messageRepo.ReplaceUserState(id, action.State)
	} else {
		uc.messageRepo.DeleteUserState(id)
	}

	return action.Text
}

func (uc *UseCase) HandleMessage(messageText string, chatID int64) string {
	id := strconv.Itoa(int(chatID))

	state, ok := uc.messageRepo.GetUserState(id)
	if !ok {
		// logic with questions or feedback
	}

	switch state.Name {
	case "login":
		return uc.ServiceLoginState(messageText, state.Step, id)
	default:
		return ""
	}
}
