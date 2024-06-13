package http

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log/slog"
)

func (s *Server) StartPolling(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	s.bot = bot

	// DELETE
	bot.Debug = true

	slog.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message != nil {
			// Command handler
			if update.Message.IsCommand() {
				text := s.useCase.HandleCommand(update.Message.Command(), update.Message.Chat.ID)
				s.SendMessage(update.Message.Chat.ID, text)
				continue
			}

			// TODO: Message handler

		} else if update.CallbackQuery != nil {
			// TODO: Callback handler

		}
	}

	return nil
}

func (s *Server) SendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)

	_, err := s.bot.Send(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}
