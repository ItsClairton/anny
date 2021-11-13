package image

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &base.Category{
	Name:         "Imagens",
	Emote:        emojis.KannaPeer,
	Interactions: []*base.Interaction{&CatCommand},
}
