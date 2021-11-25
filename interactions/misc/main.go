package misc

import "github.com/ItsClairton/Anny/core"

var Category = &core.Category{
	Name:         "Miscelâneas",
	Interactions: []*core.Interaction{&PingCommand},
}
