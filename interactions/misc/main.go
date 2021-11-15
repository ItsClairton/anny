package misc

import (
	"github.com/ItsClairton/Anny/base"
)

var Category = &base.Category{
	Name:         "Miscelâneas",
	Interactions: []*base.Interaction{&PingCommand},
}
