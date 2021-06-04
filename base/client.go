package base

import "github.com/bwmarrin/discordgo"

var (
	Client     *discordgo.Session
	commandMap = map[string]*Command{}
)

func Init(token string) error {

	var err error
	Client, err = discordgo.New("Bot " + token)

	Client.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	if err != nil {
		return err
	}

	return nil
}

func Connect() error {
	return Client.Open()
}

func AddHandler(handler interface{}) func() {
	return Client.AddHandler(handler)
}

func Disconnect() {
	Client.Close()
}

func AddCommand(cmd *Command) {

	commandMap[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		commandMap[alias] = cmd
	}

}

func GetCommandMapper() map[string]*Command {
	return commandMap
}
