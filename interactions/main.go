package interactions

import (
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/interactions/image"
	"github.com/ItsClairton/Anny/interactions/misc"
	"github.com/ItsClairton/Anny/interactions/music"
)

func init() {
	core.AddCategory(misc.Category)
	core.AddCategory(image.Category)
	core.AddCategory(music.Category)
}
