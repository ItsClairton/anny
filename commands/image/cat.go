package image

import (
	"math/rand"

	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/providers"
	"github.com/bwmarrin/discordgo"
)

var CatCommand = discord.Command{
	Name:        "cat",
	Description: "Imagem aleatória de um Gatinho",
	Options: []*discordgo.ApplicationCommandOption{{
		Name:        "gif",
		Description: "Filtrar apenas por GIF's",
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Required:    false,
	}},
	Handler: func(ctx *discord.CommandContext) {
		gif := rand.Float32() < 0.5
		if len(ctx.ApplicationCommandData().Options) > 0 {
			gif = ctx.ApplicationCommandData().Options[0].BoolValue()
		}

		info, err := providers.GetRandomCat(gif)

		if err == nil {
			ctx.SendRAW(info)
		} else {
			ctx.SendRAW("DEU ERRU " + err.Error())
		}

	},
}
