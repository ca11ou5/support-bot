package usecase

import (
	"github.com/ca11ou5/support-bot/internal/domain/message/repository/memory"
	"strconv"
	"strings"
)

func (uc *UseCase) ServiceLoginState(text string, step int, id string) (string, []memory.QA) {
	switch step {
	case 0:
		credentials := strings.Split(text, " ")
		if len(credentials) != 2 {
			return "Формат введенного сообщения не подходит под формат данных авторизации", nil
		}

		count, err := uc.messageRepo.GetSupportEmployee(credentials[0], credentials[1])
		if err != nil {
			return "Произошла внутренняя ошибка", nil
		}

		var isExist bool
		if count == 1 {
			isExist = true
		}

		if isExist {
			uc.messageRepo.IncreaseStateStep(id)
			return "Вы успешно авторизованы", []memory.QA{{
				Hash:     "hashForCreateKeyword",
				Question: "Просмотреть ключевые фразы/теги",
				// TODO
				//Answer:   "",
			}}
		}

		return "Неправильный логин или пароль", nil
	case 1:
		uc.messageRepo.DeleteUserState(id)

		chatID, _ := strconv.Atoi(id)
		uc.HandleMessage(text, int64(chatID))
	}

	return "", nil
}
