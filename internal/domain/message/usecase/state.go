package usecase

import (
	"github.com/ca11ou5/support-bot/internal/domain/message/repository/memory"
	"strconv"
	"strings"
)

func (uc *UseCase) ServiceLoginState(text string, step int, id string) (string, string, []memory.QA) {
	switch step {
	case 0:
		credentials := strings.Split(text, " ")
		if len(credentials) != 2 {
			return "Формат введенного сообщения не подходит под формат данных авторизации", "", nil
		}

		count, err := uc.messageRepo.GetSupportEmployee(credentials[0], credentials[1])
		if err != nil {
			return "Произошла внутренняя ошибка", "", nil
		}

		var isExist bool
		if count == 1 {
			isExist = true
		}

		if isExist {
			uc.messageRepo.IncreaseStateStep(id)
			return "Вы успешно авторизованы", "", []memory.QA{{
				Hash:     "hashForSeeKeyword",
				Question: "Просмотреть ключевые фразы/теги",
			},
				{
					Hash:     "hashForHelpUsers",
					Question: "Старт диалога с пользователем",
				}}
		}

		return "Неправильный логин или пароль", "", nil
	case 1:
		uc.messageRepo.DeleteUserState(id)

		chatID, _ := strconv.Atoi(id)
		uc.HandleMessage(text, int64(chatID))
	}

	return "", "", nil
}

func (uc *UseCase) ServiceChatState(text string, step int, id string) (string, string, []memory.QA) {
	opponent := uc.messageRepo.GetChatOpponent(id)

	return text, opponent, nil
}

func (uc *UseCase) ServiceAddingState(text string, step int, id string) (string, string, []memory.QA) {
	switch step {
	case 0:
		if len(text) > 64 {
			return "Размер превышает положенный", "", nil
		}

		uc.messageRepo.IncreaseStateStep(id)
		return "Введите слово еще раз для подтверждения", "", nil
	case 1:
		uc.messageRepo.SetKeyword(text)
		uc.messageRepo.DeleteUserState(id)
		return "Фраза успешно добавлена", "", nil
	}

	return "", "", nil
}
