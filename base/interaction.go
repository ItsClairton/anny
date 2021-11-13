package base

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
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
