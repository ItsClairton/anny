package misc

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &discord.Category{
	Name:     "Miscelâneas",
	Emote:    emojis.PEPE_ART,
	Commands: []*discord.Command{&PingCommand, &CatCommand},
}
