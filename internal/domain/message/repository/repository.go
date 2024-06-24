package repository

import (
	"github.com/ca11ou5/support-bot/config"
	"github.com/ca11ou5/support-bot/internal/domain/message/entity"
	"github.com/ca11ou5/support-bot/internal/domain/message/repository/memory"
	"github.com/ca11ou5/support-bot/internal/domain/message/repository/mongo"
	"github.com/ca11ou5/support-bot/internal/domain/message/repository/postgres"
	"github.com/patrickmn/go-cache"
)

type MessageRepository struct {
	mongoClient *mongo.Client
	dbClient    *postgres.Client
	memClient   *memory.Client
}

func NewMessageRepository(cfg *config.Config) *MessageRepository {
	return &MessageRepository{
		dbClient:    postgres.NewClient(cfg.PostgresURL),
		memClient:   memory.NewClient(),
		mongoClient: mongo.NewClient(cfg.MongoURL),
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

func (r *MessageRepository) GetSupportEmployee(username, password string) (int, error) {
	return r.dbClient.GetSupportEmployee(username, password)
}

func (r *MessageRepository) SaveStats(stats entity.Stats) error {
	return r.mongoClient.SaveStats(stats)
}

func (r *MessageRepository) GetStats() []entity.Stats {
	return r.mongoClient.GetStats()
}

func (r *MessageRepository) InsertWords(words map[string]int) error {
	return r.mongoClient.InsertWords(words)
}

func (r *MessageRepository) FindKeyword(words []string) []memory.QA {
	return r.memClient.FindKeyword(words)
}

func (r *MessageRepository) SetKeyword(word string) {
	r.memClient.SetKeyword(word)
}

func (r *MessageRepository) GetKeywords() map[string]cache.Item {
	return r.memClient.GetKeywords()
}

func (r *MessageRepository) AddUserToWaitChat(id string, messageText string) {
	r.memClient.AddUserToWaitChat(id, messageText)
}

func (r *MessageRepository) GetUserFromWaitList() (string, string) {
	return r.memClient.GetUserFromWaitList()
}

func (r *MessageRepository) SetupChat(userID string, id string) {
	r.memClient.SetupChat(userID, id)
}

func (r *MessageRepository) GetChatOpponent(userID string) string {
	return r.memClient.GetChatOpponent(userID)
}

func (r *MessageRepository) GetWords() map[string]interface{} {
	return r.mongoClient.GetWords()
}

func (r *MessageRepository) FindInKeywords(word string) string {
	return r.memClient.FindInKeywords(word)
}
