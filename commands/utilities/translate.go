package utilities

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/constants"
	"github.com/ItsClairton/Anny/utils/i18n"
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

		result, err := i18n.FromGoogle("auto", ctx.Args[0], strings.Join(lineArgs[1:], " "))

		if err != nil {
			ctx.ReplyWithError(err)
		} else {
			if len(result) > 2000 {
				index := strings.LastIndex(result[:2000], "\n") // Prefer lines

				if index == -1 {
					index = strings.LastIndex(result[:2000], " ")
				}

				if index == -1 {
					index = 2000
				}

				ctx.ReplyRawWithEmote(constants.PEPEFROG, result[:index])
				ctx.SendRaw(result[index:])
			} else {
				ctx.ReplyRawWithEmote(constants.PEPEFROG, result)
			}
		}
	},
}
