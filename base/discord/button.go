package discord

import (
	"math/rand"

	"github.com/ItsClairton/Anny/utils"
	"github.com/bwmarrin/discordgo"
)

var buttons = map[string]*Button{}

type Button struct {
	Label, URL, ID, UserID string

	Once    bool
	Delayed bool
	Style   discordgo.ButtonStyle
	Emoji   string
	OnClick InteractionHandler
}

func (b *Button) Build() discordgo.Button {
	RegisterButton(b)
	return discordgo.Button{
		Label:    b.Label,
		URL:      b.URL,
		Style:    b.Style,
		CustomID: b.ID,
		Emoji: discordgo.ComponentEmoji{
			Name: b.Emoji,
		},
	}
}

func GetButton(id string) *Button {
	return buttons[id]
}

func RegisterButton(b *Button) *Button {
	id := utils.Fmt("%v", rand.Float32())
	if GetButton(id) != nil {
		return RegisterButton(b)
	}
	b.ID = id
	buttons[id] = b

	return b
}

func UnregisterButton(id string) {
	buttons[id] = nil
}
