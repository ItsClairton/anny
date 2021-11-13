package misc

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &base.Category{
	Name:         "Miscelâneas",
	Emote:        emojis.PepeArt,
	Interactions: []*base.Interaction{&PingCommand},
}
