package image

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/providers"
	"github.com/diamondburned/arikawa/v3/discord"
)

var CatCommand = base.Interaction{
	Name:        "cat",
	Description: "Imagem aleatória de um Gatinho",
	Options: discord.CommandOptions{&discord.BooleanOption{
		OptionName:  "GIF",
		Description: "Filtrar apenas por GIF's",
	}},
	Handler: func(ctx *base.InteractionContext) error {
		gif := ctx.ArgumentAsBool(0)

		info, err := providers.GetRandomCat(gif)
		if err != nil {
			return ctx.SendError(err)
		}
		return ctx.Send(info)
	},
}
