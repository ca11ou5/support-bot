package usecase

import (
	"strconv"
	"strings"
)

func (uc *UseCase) ServiceLoginState(text string, step int, id string) string {
	switch step {
	case 0:
		credentials := strings.Split(text, " ")
		if len(credentials) != 2 {
			return "Формат введенного сообщения не подходит под формат данных авторизации"
		}

		count, err := uc.messageRepo.GetSupportEmployee(credentials[0], credentials[1])
		if err != nil {
			return "Произошла внутренняя ошибка"
		}

		var isExist bool
		if count == 1 {
			isExist = true
		}

		if isExist {
			uc.messageRepo.IncreaseStateStep(id)
			return "Вы успешно авторизованы"
		}

		return "Неправильный логин или пароль"
	case 1:
		uc.messageRepo.DeleteUserState(id)

		chatID, _ := strconv.Atoi(id)
		uc.HandleMessage(text, int64(chatID))
	}

	return ""
}
