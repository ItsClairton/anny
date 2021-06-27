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

		result, err := i18n.FromGoogle("auto", ctx.Args[0], strings.Join(lineArgs[1:], " "))

		if err != nil {
			ctx.ReplyWithError(err)
		} else {
			if len(result) > 1024 {

				spaceIndex := strings.LastIndex(result[:1024], " ")
				var firstPart string
				var secondPart string

				if spaceIndex > -1 {
					firstPart = result[:spaceIndex]
					secondPart = result[spaceIndex:]
				} else {
					firstPart = result[:1024]
					secondPart = result[1024:]
				}

				ctx.ReplyRawWithEmote(constants.PEPEFROG, firstPart)
				ctx.SendRaw(secondPart)
			} else {
				ctx.ReplyRawWithEmote(constants.PEPEFROG, result)
			}
		}
	},
}
