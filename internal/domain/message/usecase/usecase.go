package usecase

import (
	"github.com/ca11ou5/support-bot/internal/domain/message/repository"
	"github.com/ca11ou5/support-bot/internal/domain/message/repository/memory"
	"strconv"
	"strings"
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

func (uc *UseCase) HandleMessage(messageText string, chatID int64) (string, []memory.QA) {
	id := strconv.Itoa(int(chatID))

	state, ok := uc.messageRepo.GetUserState(id)
	if !ok {
		words := strings.Split(messageText, " ")

		qa := uc.messageRepo.FindKeyword(words)
		if len(qa) != 0 {
			return "По вашему запросу найдены похожие подсказки:", qa
		}

		return "По вашему запросу ничего не найдено, связаться со специалистом?", []memory.QA{
			{
				Hash:     "hashForConfirm",
				Question: "Да",
			}, {
				Hash:     "hashForDeny",
				Question: "Нет",
			},
		}
	}

	switch state.Name {
	case "login":
		return uc.ServiceLoginState(messageText, state.Step, id)
	default:
		return "", nil
	}
}

func (uc *UseCase) HandleCallback(callbackData string, chatID int64) CallbackAnswer {
	switch callbackData {
	case "hashForSeeKeyword":
		return uc.HandleSeeKeyword()
	default:
		return CallbackAnswer{}
	}
}
