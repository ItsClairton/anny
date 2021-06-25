package utilities

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/constants"
)

var Category = &base.Category{
	ID:       "utilities",
	Emote:    constants.PEPEPOGGERS,
	Commands: []*base.Command{&AnimeCommand, &MangaCommand, &SceneCommand},
}
