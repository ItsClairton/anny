package miscellaneous

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/constants"
)

var Category = &base.Category{
	ID:       "miscellaneous",
	Emote:    constants.ZERO_HMPF,
	Commands: []*base.Command{&PingCommand},
}
