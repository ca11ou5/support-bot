package http

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (h *Handler) ListenAndServeBotAPI(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			DispatchMessage(bot, update)
		} else if update.CallbackQuery != nil {
			DispatchCallback(bot, update)
		}
	}
}
