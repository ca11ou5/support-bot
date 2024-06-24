package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
)

type CallbackAnswer struct {
	Text              string
	MessageKeyboard   tgbotapi.InlineKeyboardMarkup
	OpponentID        string
	MessageToOpponent string
	KB                tgbotapi.ReplyKeyboardMarkup
}

func (uc *UseCase) HandleSeeKeyword() CallbackAnswer {
	keywords := uc.messageRepo.GetKeywords()
	if len(keywords) == 0 {
		return CallbackAnswer{
			Text: "В данный момент не существует ни одного ключевого слова/фразы",
			MessageKeyboard: tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить ключевое"+
				" слово/фразу", "AddNewKeyword"))),
		}
	}

	ca := CallbackAnswer{
		Text:            "Список всех ключевых слов и фраз",
		MessageKeyboard: generateCallbackKeyboard(keywords),
	}

	return ca
}

func generateCallbackKeyboard(keywords map[string]cache.Item) tgbotapi.InlineKeyboardMarkup {
	var keyboard tgbotapi.InlineKeyboardMarkup

	for i, _ := range keywords {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(i, i)})
	}

	return keyboard
}

func (uc *UseCase) HandleConfirm(id string, messageText string) CallbackAnswer {
	uc.messageRepo.AddUserToWaitChat(id, messageText)

	return CallbackAnswer{
		Text: "Ваш вопрос отправлен специалисту, ожидайте ответа",
	}
}

func (uc *UseCase) HandleHelpUsers(id string) CallbackAnswer {
	userID, messageText := uc.messageRepo.GetUserFromWaitList()
	if userID == "" {
		return CallbackAnswer{Text: "В данный момент никаких вопросов нет, хорошая работа! =)"}
	}

	uc.messageRepo.ReplaceUserState(id, "chat")
	uc.messageRepo.ReplaceUserState(userID, "chat")

	uc.messageRepo.SetupChat(userID, id)
	return CallbackAnswer{
		Text:              "Постарайтесь дать наиболее развернутый ответ пользователю\n\nВопрос пользователя:\n" + messageText,
		OpponentID:        userID,
		MessageToOpponent: "На связи специалист Александр",
		KB: tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Закончить диалог"),
			)),
	}
}

func (uc *UseCase) HandleAddKeyword(id string) CallbackAnswer {
	uc.messageRepo.ReplaceUserState(id, "adding")
	return CallbackAnswer{
		Text: "Введите ключевое слово/фразу\nДлина фразы не должна превышать 64 символа",
	}
}

func (uc *UseCase) HandleAddQA(id string, word string) CallbackAnswer {
	uc.messageRepo.ReplaceUserState(id, "addingQA")
	return CallbackAnswer{
		Text:              "",
		MessageKeyboard:   tgbotapi.InlineKeyboardMarkup{},
		OpponentID:        "",
		MessageToOpponent: "",
		KB:                tgbotapi.ReplyKeyboardMarkup{},
	}
}
