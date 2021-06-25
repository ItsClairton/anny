package miscellaneous

import (
	"os"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/constants"
)

var HelpCommand = base.Command{
	Name:    "help",
	Aliases: []string{"ajuda", "h"},
	Handler: func(ctx *base.CommandContext) {
		ch, err := ctx.Client.UserChannelCreate(ctx.Author.ID)

		if err != nil {
			ctx.ReplyWithError(err)
			return
		}

		eb := embed.NewEmbed(ctx.Locale, "miscellaneous.help.embed").
			WithDescription(ctx.Author.Username).
			SetColor(0x7289DA)

		for _, category := range base.GetCategories() {
			var commands string

			for _, cmd := range category.Commands {
				commands += utils.Fmt("`%s%s` - %s\n", os.Getenv("DEFAULT_PREFIX"), cmd.Name, ctx.GetString(utils.Fmt("%s.%s.description", category.ID, cmd.Name)))
			}

			eb.AddField(utils.Fmt("%s %s (%d)", category.Emote, ctx.GetString(utils.Fmt("%s.categoryName", category.ID)), len(category.Commands)), commands, false)
		}

		_, err = ctx.SendWithEmbedTo(ch.ID, eb)

		if err != nil {
			ctx.Reply(constants.MIKU_CRY, "miscellaneous.help.dmClosed")
		} else {
			ctx.Reply(constants.PEPEFROG, "miscellaneous.help.reply")
		}
	},
}
