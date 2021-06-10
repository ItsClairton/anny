package image

import (
	"math/rand"
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/services/image"
)

var NekoCommand = base.Command{
	Name: "neko",
	Handler: func(ctx *base.CommandContext) {

		gif := ctx.Args != nil && strings.HasPrefix(ctx.Args[0], "gif") || rand.Float32() < 0.2

		var url string
		var err error

		if gif {
			url, err = image.GetFromNekos("ngif")
		} else {
			if rand.Float32() < 0.5 {
				url, err = image.GetFromNekoBot("neko")
			} else {
				url, err = image.GetFromNekos("neko")
			}
		}

		if err != nil {
			ctx.ReplyWithError(err)
		} else {
			ctx.ReplyRaw(url)
		}
	},
}
