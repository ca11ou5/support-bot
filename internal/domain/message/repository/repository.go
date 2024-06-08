package repository

import "github.com/ca11ou5/support-bot/internal/domain/message/repository/memory"

type MessageRepository struct {
	// MongoClient
	// PostgresClient
	memClient *memory.Client
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		memClient: memory.NewClient(),
	}
}

func (r *MessageRepository) GetCommandAction(command string) *memory.CommandAction {
	return r.memClient.GetCommandAction(command)
}
