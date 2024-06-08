package memory

type CommandAction struct {
	Text             string
	NeedsStateChange bool
}

var commandPhrases = map[string]CommandAction{
	"start": {
		Text:             "Вы вели стартовую команду",
		NeedsStateChange: false,
	},
}

type Client struct {
	commandPhrases map[string]CommandAction
	// stateMachine
}

func NewClient() *Client {
	return &Client{
		commandPhrases: commandPhrases,
	}
}

func (c *Client) GetCommandAction(command string) *CommandAction {
	v, ok := c.commandPhrases[command]
	if !ok {
		return &CommandAction{
			Text:             "Я не знаю такого =(",
			NeedsStateChange: false,
		}
	}

	return &v
}
