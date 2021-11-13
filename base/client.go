package base

import (
	"context"

	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	Session      *state.State
	Interactions = map[string]*Interaction{}
	categories   = []*Category{}
)

func New(token string) (err error) {
	Session, err = state.NewWithIntents(utils.Fmt("Bot %s", token), gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentGuildVoiceStates)

	if err == nil {
		err = Session.Open(context.Background())
	}

	return err
}

func CheckInteractions() error {
	app, err := Session.CurrentApplication()
	if err != nil {
		return err
	}

	previous, err := Session.Commands(app.ID)
	if err != nil {
		return err
	}

	checked := map[string]api.CreateCommandData{}
	for _, prevIn := range previous {
		newIn := Interactions[prevIn.Name]

		if newIn == nil { // Remover interações antigas que não existem mais no bot.
			logger.DebugF("Removendo interação \"%s\" do Discord...", prevIn.Name)

			if err := Session.DeleteCommand(app.ID, prevIn.ID); err != nil {
				logger.WarnF("Não foi possível remover a interação \"%s\" do Discord: %v", prevIn.Name, err)
			}
		} else { // Atualizar informações da interação no Discord, caso elas não estejam atualizadas.
			// TODO: DeepEqual não funciona de forma correta com Arikawa
			if len(prevIn.Options) != len(newIn.Options) || newIn.Description != prevIn.Description {
				logger.DebugF("Atualizando interação \"%s\" no Discord...", newIn.Name)

				if _, err := Session.EditCommand(app.ID, prevIn.ID, newIn.RAW()); err != nil {
					logger.WarnF("Não foi possivel atualizar as informações da interação \"%s\" no Discord: %v", newIn.Name, err)
				}
			}

			checked[newIn.Name] = newIn.RAW()
		}
	}

	for _, interaction := range Interactions {
		if _, exist := checked[interaction.Name]; !exist {
			logger.DebugF("Criando interação \"%s\" no Discord...", interaction.Name)

			if _, err := Session.CreateCommand(app.ID, interaction.RAW()); err != nil {
				logger.WarnF("Não foi possivel criar interação \"%s\" no Discord: %v", interaction.Name, err)
			}
		}
	}

	return nil
}

func AddHandler(handler interface{}) {
	Session.AddHandler(handler)
}

func AddCategory(category *Category) {
	for _, interaction := range category.Interactions {
		interaction.Category = category
		Interactions[interaction.Name] = interaction
	}

	categories = append(categories, category)
}

func SendMessage(id discord.ChannelID, emoji, text string, args ...interface{}) {
	Session.SendMessage(id, utils.Fmt("%s | %s", emoji, utils.Fmt(text, args...)))
}

func Disconnect() {
	Session.Close()
}
