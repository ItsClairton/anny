package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Category struct {
	Name, Emote  string
	Interactions []*Interaction
}

type Interaction struct {
	Name, Description string
	Type              discordgo.ApplicationCommandType
	Deffered          bool
	Options           []*discordgo.ApplicationCommandOption
	Category          *Category
	Handler           InteractionHandler
}

type InteractionHandler func(*InteractionContext)

func (i Interaction) ToRAW() *discordgo.ApplicationCommand {
	raw := &discordgo.ApplicationCommand{
		Name:        i.Name,
		Description: i.Description,
		Type:        i.Type,
		Options:     i.Options,
	}

	return raw
}
