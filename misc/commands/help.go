package misc

import (
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils"
)

var HelpCommand = core.Command{
	Name:        "ajuda",
	Description: "Pagina de comandos",
	Handler: func(ctx *core.CommandContext) {
		embed := utils.NewEmbed().Color(0x7289da)

		for _, module := range core.Modules {
			if len(module.Commands) == 0 {
				continue
			}

			var fieldDesc string
			for _, cmd := range module.Commands {
				fieldDesc += utils.Fmt("`/%s` - %s\n", cmd.Name, cmd.Description)
			}

			embed.Field(utils.Fmt("%s %s", module.Emote, module.Name), fieldDesc, false)
		}

		ctx.Ephemeral().Reply(embed)
	},
}
