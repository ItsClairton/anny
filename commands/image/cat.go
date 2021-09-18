package image

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/providers/image"
)

var CatCommand = base.Command{
	Name:    "cat",
	Aliases: []string{"gato", "meow"},
	Handler: func(ctx *base.CommandContext) {
		url, err := image.GetRandomCat(ctx.Args != nil && strings.HasPrefix(ctx.Args[0], "gif"))

		if err != nil {
			ctx.ReplyWithError(err)
		} else {
			ctx.ReplyRaw(url)
		}
	},
}
