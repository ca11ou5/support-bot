package main

import (
	"github.com/ca11ou5/support-bot/config"
	"github.com/ca11ou5/support-bot/internal/controller/http"
	"github.com/ca11ou5/support-bot/internal/domain/message/repository"
	"github.com/ca11ou5/support-bot/internal/domain/message/usecase"
)

func main() {
	var cfg config.Config

	// TODO: config reading

	repo := repository.NewMessageRepository()
	uc := usecase.NewMessageUseCase(repo)
	server := http.NewServer(uc)

	err := server.Start(&cfg)
	if err != nil {
		return
	}
}
