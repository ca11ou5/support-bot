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

//func (r *MessageRepository) ChangeUserState(userID, action string) {
//	r.memClient.ChangeUserState(userID, action)
//}

func (r *MessageRepository) DeleteUserState(chatID string) {
	r.memClient.DeleteUserState(chatID)
}

func (r *MessageRepository) ReplaceUserState(chatID string, newState string) {
	r.memClient.ReplaceUserState(chatID, newState)
}

func (r *MessageRepository) GetUserState(chatID string) (memory.State, bool) {
	return r.memClient.GetUserState(chatID)
}

func (r *MessageRepository) IncreaseStateStep(chatID string) {
	r.memClient.IncreaseStateStep(chatID)
}
