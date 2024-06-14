package usecase

import "strings"

func (uc *UseCase) ServiceLoginState(text string, step int, id string) string {
	switch step {
	case 0:
		credentials := strings.Split(text, " ")
		if len(credentials) != 2 {
			return "Формат введенного сообщения не подходит под формат данных авторизации"
		}

		// TODO: check credentials in sql database
		//if Valid {
		uc.messageRepo.IncreaseStateStep(id)
		return "Вы успешно авторизованы"
		//}

		return "Неправильный логин или пароль"
	}
}
