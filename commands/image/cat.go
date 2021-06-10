package image

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/services/image"
)

var CatCommand = base.Command{
	Name:    "gato",
	Aliases: []string{"cat", "meow"},
	Handler: func(ctx *base.CommandContext) {
		url, err := image.GetRandomCat(ctx.Args != nil && strings.HasPrefix(ctx.Args[0], "gif"))

		if err != nil {
			ctx.ReplyWithError(err)
		} else {
			ctx.ReplyRaw(url)
		}
	},
}
