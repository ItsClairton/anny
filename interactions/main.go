package interactions

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/interactions/image"
	"github.com/ItsClairton/Anny/interactions/misc"
	"github.com/ItsClairton/Anny/interactions/music"
)

func init() {
	base.AddCategory(misc.Category)
	base.AddCategory(image.Category)
	base.AddCategory(music.Category)
}
