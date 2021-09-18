package discord

import (
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	Session    *discordgo.Session
	commands   = map[string]*Command{}
	categories = []*Category{}
)

func Init(token string) {
	Session, _ = discordgo.New("Bot " + token)

	Session.Identify.Intents = discordgo.IntentsAll
}

func Connect() error {
	return Session.Open()
}

func Disconnect() {
	Session.Close()
}

func AddCategory(category *Category) {
	for _, cmd := range category.Commands {
		err := addCommand(cmd, category)

		if err != nil {
			logger.Error("[%s] Um erro ocorreu ao registrar o comando %s no Discord. (%s)", category.Name, cmd.Name, err.Error())
		}
	}

	categories = append(categories, category)
}

func addCommand(cmd *Command, category *Category) error {
	cmd.Category = category

	data := &discordgo.ApplicationCommand{
		Name:        cmd.Name,
		Description: cmd.Description,
		Type:        cmd.Type,
	}
	if cmd.Options != nil {
		data.Options = cmd.Options
	}

	_, err := Session.ApplicationCommandCreate(Session.State.User.ID, "", data)

	if err == nil {
		commands[cmd.Name] = cmd
	}

	return err
}

func GetCommands() map[string]*Command {
	return commands
}
