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
		lineArgs := ctx.GetArgsWithLines()

		if len(lineArgs) < 2 {
			ctx.ReplyWithUsage("<linguagem> <texto>")
			return
		}

		result, err := i18n.FromGoogle("auto", lineArgs[0], strings.Join(lineArgs[1:], " "))

		if err != nil {
			ctx.ReplyWithError(err)
		} else {
			ctx.ReplyRawWithEmote(constants.PEPEFROG, result)
		}
	},
}
