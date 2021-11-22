package base

import (
	"reflect"

	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

var (
	Interactions = map[string]*Interaction{}
	categories   = []*Category{}
)

type Category struct {
	Name, Emote  string
	Interactions []*Interaction
}

type Interaction struct {
	Name, Description string
	Type              discord.CommandType
	Deffered          bool
	Options           discord.CommandOptions
	Category          *Category
	Handler           InteractionHandler
}

type InteractionHandler func(*InteractionContext) error

func (i Interaction) RAW() api.CreateCommandData {
	return api.CreateCommandData{
		Name:        i.Name,
		Description: i.Description,
		Type:        i.Type,
		Options:     i.Options,
	}
}

func AddCategory(category *Category) {
	for _, interaction := range category.Interactions {
		interaction.Category = category
		Interactions[interaction.Name] = interaction
	}

	categories = append(categories, category)
}

func DeployInteractions() error {
	app, err := Session.CurrentApplication()
	if err != nil {
		return err
	}

	previous, err := Session.Commands(app.ID)
	if err != nil {
		return err
	}

	checked := []string{}
	for _, prevIn := range previous {
		newIn := Interactions[prevIn.Name]

		if newIn == nil { // Remover interação antiga que não existe mais no bot.
			logger.DebugF("Removendo interação \"%s\" do Discord...", prevIn.Name)

			if err := Session.DeleteCommand(app.ID, prevIn.ID); err != nil {
				logger.WarnF("Não foi possível remover a interação \"%s\" do Discord: %v", prevIn.Name, err)
			}
		} else { // Atualizar informações da interação no Discord, caso elas não estejam atualizadas.
			if !reflect.DeepEqual(prevIn.Options, newIn.Options) || newIn.Description != prevIn.Description {
				logger.DebugF("Atualizando interação \"%s\" no Discord...", newIn.Name)

				if _, err := Session.EditCommand(app.ID, prevIn.ID, newIn.RAW()); err != nil {
					logger.WarnF("Não foi possivel atualizar as informações da interação \"%s\" no Discord: %v", newIn.Name, err)
				}
			}

			checked = append(checked, newIn.Name)
		}
	}

	for _, interaction := range Interactions {
		if !utils.StringArrayContains(checked, interaction.Name) {
			logger.DebugF("Criando interação \"%s\" no Discord...", interaction.Name)

			if _, err := Session.CreateCommand(app.ID, interaction.RAW()); err != nil {
				logger.WarnF("Não foi possivel criar interação \"%s\" no Discord: %v", interaction.Name, err)
			}
		}
	}

	return nil
}
