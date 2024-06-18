package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
)

type CallbackAnswer struct {
	Text            string
	MessageKeyboard tgbotapi.InlineKeyboardMarkup
}

func (uc *UseCase) HandleSeeKeyword() CallbackAnswer {
	keywords := uc.messageRepo.GetKeywords()
	if len(keywords) == 0 {
		return CallbackAnswer{
			Text: "В данный момент не существует ни одного ключевого слова/фразы",
		}
	}

	ca := CallbackAnswer{
		Text:            "Список всех ключевых слов",
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
