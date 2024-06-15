package main

import (
	"github.com/ca11ou5/support-bot/config"
	"github.com/ca11ou5/support-bot/internal/controller/http"
	"github.com/ca11ou5/support-bot/internal/domain/message/repository"
	"github.com/ca11ou5/support-bot/internal/domain/message/usecase"
	"github.com/ca11ou5/support-bot/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
)

func main() {
	logging.SetupLogger()

	var cfg config.Config

	err := cleanenv.ReadConfig("./envs/dev.env", &cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	repo := repository.NewMessageRepository(&cfg)
	uc := usecase.NewMessageUseCase(repo)
	server := http.NewServer(uc)

	err = server.Start(&cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}

}
