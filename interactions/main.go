package interactions

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/interactions/image"
	"github.com/ItsClairton/Anny/interactions/misc"
)

func init() {
	discord.AddCategory(misc.Category)
	discord.AddCategory(image.Category)
}
