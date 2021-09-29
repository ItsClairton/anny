package discord

import (
	"reflect"

	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	Session      *discordgo.Session
	interactions = map[string]*Interaction{}
	categories   = []*Category{}
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
	for _, i := range category.Interactions {
		addInteraction(i, category)
	}

	categories = append(categories, category)
}

func addInteraction(i *Interaction, category *Category) {
	i.Category = category
	interactions[i.Name] = i
}

func GetInteractions() map[string]*Interaction {
	return interactions
}

func RegisterInDiscord() error {
	previous, err := Session.ApplicationCommands(Session.State.User.ID, "")
	if err != nil {
		return err
	}

	registered := map[string]*discordgo.ApplicationCommand{}
	for _, prev := range previous { // Verificar se precisa atualizar ou remover alguma interação no Discord.
		i, exist := interactions[prev.Name]

		if !exist {
			err := Session.ApplicationCommandDelete(Session.State.User.ID, "", prev.ID)
			if err != nil {
				logger.Warn("Não foi possível remover a interação \"%s\" do Discord. (%s)", prev.Name, err.Error())
			}
		} else {
			if !reflect.DeepEqual(i.Options, prev.Options) || prev.Description != i.Description {
				_, err = Session.ApplicationCommandEdit(Session.State.User.ID, "", prev.ID, i.ToRAW())
				if err != nil {
					logger.Warn("Não foi possível atualizar a interação \"%s\" no Discord. (%s)", i.Name, err.Error())
				} else {
					registered[i.Name] = i.ToRAW()
				}
			} else {
				registered[i.Name] = i.ToRAW()
			}
		}
	}

	for _, i := range interactions { // Registrar novas interações no Discord caso não existam.
		_, exist := registered[i.Name]

		if !exist {
			_, err := Session.ApplicationCommandCreate(Session.State.User.ID, "", i.ToRAW())
			if err != nil {
				logger.Warn("Não foi possível criar a interação \"%s\" no Discord. (%s)", i.Name, err.Error())
			}
		}
	}

	return nil
}
