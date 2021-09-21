package image

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/providers"
	"github.com/bwmarrin/discordgo"
)

var CatCommand = discord.Interaction{
	Name:        "cat",
	Description: "Imagem aleatória de um Gatinho",
	Options: []*discordgo.ApplicationCommandOption{{
		Name:        "gif",
		Description: "Filtrar apenas por GIF's",
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Required:    false,
	}},
	Handler: func(ctx *discord.InteractionContext) {
		gif := len(ctx.ApplicationCommandData().Options) > 0 && ctx.ApplicationCommandData().Options[0].BoolValue()
		info, err := providers.GetRandomCat(gif)

		if err == nil {
			ctx.SendRAW(info)
		} else {
			ctx.SendError(err)
		}

	},
}
