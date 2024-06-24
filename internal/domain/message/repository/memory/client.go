package memory

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type CommandAction struct {
	Text            string
	NeedsSetupState bool
	State           string
	Menu            []QA
}

var commandPhrases = map[string]CommandAction{
	"start": {
		Text:            "Вы вели стартовую команду",
		NeedsSetupState: false,
	},
	"login": {
		Text:            "Введите данные для авторизации в формате:\nTODO: прикрутить Mini Apps",
		NeedsSetupState: true,
		State:           "login",
	},
	"menu": {
		Text:            "Навигационное меню",
		NeedsSetupState: false,
		Menu: []QA{{
			Hash:     "hashForSeeKeyword",
			Question: "Просмотреть ключевые фразы/теги",
		},
			{
				Hash:     "hashForHelpUsers",
				Question: "Старт диалога с пользователем",
			}},
	},
}

type Client struct {
	commandPhrases  map[string]CommandAction
	states          *cache.Cache
	keywords        *cache.Cache
	waitingForChats map[string]string
	chats           *cache.Cache
}

func NewClient() *Client {
	return &Client{
		commandPhrases:  commandPhrases,
		states:          cache.New(7*time.Hour*24, 10*time.Minute),
		keywords:        cache.New(365*time.Hour*24, 10*time.Minute),
		waitingForChats: make(map[string]string),
		chats:           cache.New(365*time.Hour*24, 10*time.Minute),
	}
}

func (c *Client) GetCommandAction(command string) *CommandAction {
	v, ok := c.commandPhrases[command]
	if !ok {
		return &CommandAction{
			Text:            "Я не знаю такого =(",
			NeedsSetupState: false,
		}
	}

	return &v
}

type State struct {
	Name string
	Step int
}

func (c *Client) DeleteUserState(userID string) {
	c.states.Delete(userID)
}

func (c *Client) ReplaceUserState(userID string, newState string) {
	c.states.Delete(userID)

	c.states.Set(userID, State{
		Name: newState,
		Step: 0,
	}, cache.DefaultExpiration)
}

func (c *Client) GetUserState(userID string) (State, bool) {
	state, ok := c.states.Get(userID)
	if !ok {
		return State{}, false
	}

	return state.(State), true
}

func (c *Client) IncreaseStateStep(userID string) {
	state, _ := c.states.Get(userID)

	newState := State{
		Name: state.(State).Name,
		Step: state.(State).Step + 1,
	}

	c.states.Set(userID, newState, cache.DefaultExpiration)
}

type Keyword struct {
	Word  string
	QA    []QA
	Count int
}

type QA struct {
	Hash     string
	Question string
	Answer   string
}

func (c *Client) FindKeyword(words []string) []QA {
	var qa []QA

	for _, word := range words {
		v, ok := c.keywords.Get(word)
		if ok {
			keyword := v.(Keyword)
			keyword.Count++

			c.keywords.Set(word, keyword, cache.DefaultExpiration)
			for _, val := range keyword.QA {
				qa = append(qa, val)
			}
		}
	}

	return qa
}

func (c *Client) FindInKeywords(word string) string {
	_, ok := c.keywords.Get(word)
	if ok {
		return word
	}

	return ""
}

func (c *Client) SetKeyword(word string) {
	kw := Keyword{
		Word: word,
		QA:   []QA{},
	}

	c.keywords.Set(word, kw, cache.DefaultExpiration)
}

func (c *Client) GetKeywords() map[string]cache.Item {
	return c.keywords.Items()
}

func (c *Client) AddUserToWaitChat(id string, messageText string) {
	c.waitingForChats[id] = messageText
}

func (c *Client) GetUserFromWaitList() (string, string) {
	for k, v := range c.waitingForChats {
		delete(c.waitingForChats, k)
		return k, v
	}

	return "", ""
}

type Chat struct {
	UserID     string
	EmployerID string
}

func (c *Client) SetupChat(userID string, id string) {
	chat := Chat{
		UserID:     userID,
		EmployerID: id,
	}

	c.chats.Set(userID, chat, cache.DefaultExpiration)
	c.chats.Set(id, chat, cache.DefaultExpiration)
}

func (c *Client) GetChatOpponent(id string) string {
	chat, ok := c.chats.Get(id)
	if !ok {
		return ""
	}

	chatT := chat.(Chat)
	if chatT.EmployerID != id {
		return chatT.EmployerID
	}

	return chatT.UserID
}
