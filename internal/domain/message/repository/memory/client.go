package memory

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type CommandAction struct {
	Text            string
	NeedsSetupState bool
	State           string
}

var commandPhrases = map[string]CommandAction{
	"start": {
		Text:            "Вы вели стартовую команду",
		NeedsSetupState: false,
	},
}

type Client struct {
	commandPhrases map[string]CommandAction
	states         *cache.Cache
}

func NewClient() *Client {
	return &Client{
		commandPhrases: commandPhrases,
		states:         cache.New(7*time.Hour*24, 10*time.Minute),
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