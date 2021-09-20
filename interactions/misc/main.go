package misc

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &discord.Category{
	Name:         "Miscelâneas",
	Emote:        emojis.PepeArt,
	Interactions: []*discord.Interaction{&PingCommand},
}
