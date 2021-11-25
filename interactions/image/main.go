package image

import (
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &core.Category{
	Name:         "Imagens",
	Emote:        emojis.Peer,
	Interactions: []*core.Interaction{&CatCommand},
}
