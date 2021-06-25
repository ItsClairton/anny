package base

import (
	"github.com/ItsClairton/Anny/i18n"
	"github.com/ItsClairton/Anny/logger"
	"github.com/ItsClairton/Anny/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	Client     *discordgo.Session
	commandMap = map[string]*Command{}
	categories = []*Category{}
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

func AddCategory(category *Category) {
	if i18n.GetDefaultLocale().GetString(utils.Fmt("%s.categoryName", category.ID)) == "N/A" {
		logger.Warn("Não foi possível encontrar o nome da categoria com o ID %s no arquivo de tradução padrão, portanto ela não será carregada.", category.ID)
	} else {
		for _, i := range category.Commands {
			addCommand(i, category)
		}

		categories = append(categories, category)
	}
}

func addCommand(cmd *Command, category *Category) {
	if i18n.GetDefaultLocale().GetString(utils.Fmt("%s.%s.description", category.ID, cmd.Name)) == "N/A" {
		logger.Warn("Não foi possível encontrar a descrição do comando %s no arquivo de tradução padrão, portanto ele não será carregado.", cmd.Name)
	} else {
		cmd.Category = category

		commandMap[cmd.Name] = cmd
		for _, alias := range cmd.Aliases {
			commandMap[alias] = cmd
		}
	}
}

func GetCommandMapper() map[string]*Command {
	return commandMap
}

func GetCategories() []*Category {
	return categories
}
