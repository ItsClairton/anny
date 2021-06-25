package utilities

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/i18n"
	"github.com/ItsClairton/Anny/utils/constants"
)

var TranslateCommand = base.Command{
	Name:    "translate",
	Aliases: []string{"traduzir"},
	Handler: func(ctx *base.CommandContext) {
		if len(ctx.Args) < 2 {
			ctx.ReplyWithUsage("<linguagem> <texto>")
			return
		}

		result, err := i18n.FromGoogle("auto", ctx.Args[0], strings.Join(ctx.Args[1:], " "))

		if err != nil {
			ctx.ReplyWithError(err)
		} else {
			ctx.ReplyRawWithEmote(constants.PEPEFROG, result)
		}
	},
}
