package http

import (
	"github.com/ca11ou5/support-bot/internal/domain/message/repository/memory"
	"github.com/ca11ou5/support-bot/internal/domain/message/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log/slog"
	"strconv"
)

func (s *Server) StartPolling(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	s.bot = bot

	// DELETE
	//bot.Debug = true
	bot.Debug = false

	slog.Info("Authorized on account", "username", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	go s.HandleStats()

	for update := range updates {
		go s.statsCounting(update)

		if update.Message != nil {
			// Command handler
			if update.Message.IsCommand() {
				text := s.useCase.HandleCommand(update.Message.Command(), update.Message.Chat.ID)
				s.SendMessage(update.Message.Chat.ID, text, false)
				continue
			}

			text, id, qas := s.useCase.HandleMessage(update.Message.Text, update.Message.Chat.ID)
			if id != "" {
				ids, _ := strconv.Atoi(id)

				if text == "Диалог закончен" {
					s.SendMessage(int64(ids), text, true)
					s.SendMessage(update.Message.Chat.ID, text, true)
					continue
				}

				s.SendMessage(int64(ids), text, false)
				continue
			}

			if qas != nil {
				s.SendKeyboard(update.Message.Chat.ID, text, qas, update.Message.MessageID)
				continue
			}

			s.SendMessage(update.Message.Chat.ID, text, false)
			continue

		} else if update.CallbackQuery != nil {
			ca := s.useCase.HandleCallback(update.CallbackQuery.Data, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.ReplyToMessage.Text)
			s.SendCA(ca, update.CallbackQuery.Message.Chat.ID)
			continue
		}
	}

	return nil
}

func (s *Server) SendMessage(chatID int64, message string, needCloseKB bool) {
	msg := tgbotapi.NewMessage(chatID, message)
	if needCloseKB {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}

	_, err := s.bot.Send(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (s *Server) SendKeyboard(chatID int64, message string, qas []memory.QA, messageID int) {
	msg := tgbotapi.NewMessage(chatID, message)
	var keyboard tgbotapi.InlineKeyboardMarkup

	for _, qa := range qas {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(qa.Question, qa.Hash)})
	}

	msg.ReplyMarkup = keyboard
	msg.ReplyToMessageID = messageID

	_, err := s.bot.Send(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (s *Server) SendCA(ca usecase.CallbackAnswer, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, ca.Text)

	if len(ca.MessageKeyboard.InlineKeyboard) != 0 {
		msg.ReplyMarkup = ca.MessageKeyboard
	}

	if ca.OpponentID != "" {
		opID, _ := strconv.Atoi(ca.OpponentID)
		secMsg := tgbotapi.NewMessage(int64(opID), ca.MessageToOpponent)

		msg.ReplyMarkup = ca.KB
		secMsg.ReplyMarkup = ca.KB
		s.bot.Send(secMsg)
	}

	_, err := s.bot.Send(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}
