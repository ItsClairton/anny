package discord

import (
	"reflect"

	"github.com/ItsClairton/Anny/utils"
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

func SendMessage(ChannelID, emoji, content string, args ...interface{}) (*discordgo.Message, error) {
	return Session.ChannelMessageSend(ChannelID, utils.Fmt("%s | %s", emoji, utils.Fmt(content, args...)))
}

func GetInteractions() map[string]*Interaction {
	return interactions
}

func UpdateInteractions() error {
	logger.Debug("Verificando se há interações para atualizar ou remover do Discord...")

	previous, err := Session.ApplicationCommands(Session.State.User.ID, "")
	if err != nil {
		return err
	}

	checked := map[string]*discordgo.ApplicationCommand{}
	for _, prevIn := range previous { // Atualizar ou remover interações, caso necessário
		newIn := interactions[prevIn.Name]
		if newIn == nil { // Remover interação, caso ela não exista no mapeador de interações.
			logger.Debug(utils.Fmt("Deletando interação \"%s\" do Discord...", prevIn.Name))

			if err := Session.ApplicationCommandDelete(Session.State.User.ID, "", prevIn.ID); err != nil {
				logger.Warn(utils.Fmt("Não foi possível remover a interação \"%s\" do Discord.", prevIn.Name), err)
			}
		} else { // Atualizar interação, caso a descrição, ou as opções da interação estejam desatualizadas no Discord.
			if !reflect.DeepEqual(newIn.Options, prevIn.Options) || prevIn.Description != newIn.Description {
				logger.Debug(utils.Fmt("Atualizando interação \"%s\" no Discord...", newIn.Name))

				if _, err := Session.ApplicationCommandEdit(Session.State.User.ID, "", prevIn.ID, newIn.ToRAW()); err != nil {
					logger.Warn(utils.Fmt("Não foi possível atualizar a interação \"%s\" no Discord.", newIn.Name), err)
				}
			}
			checked[newIn.Name] = newIn.ToRAW()
		}
	}

	for _, newIn := range interactions { // Criar novas interações, se necessário
		if checked[newIn.Name] == nil {
			logger.Debug(utils.Fmt("Criando interação \"%s\" no Discord...", newIn.Name))

			if _, err := Session.ApplicationCommandCreate(Session.State.User.ID, "", newIn.ToRAW()); err != nil {
				logger.Warn(utils.Fmt("Não foi possível criar a interação \"%s\" no Discord.", newIn.Name), err)
			}
		}
	}

	return nil
}

func addInteraction(i *Interaction, category *Category) {
	i.Category = category
	interactions[i.Name] = i
}
