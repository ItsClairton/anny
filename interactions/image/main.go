package image

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &discord.Category{
	Name:         "Imagens",
	Emote:        emojis.KannaPeer,
	Interactions: []*discord.Interaction{&CatCommand, &TraceContext},
}