package image

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/services/ImageService"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/sutils"
)

var CatCommand = base.Command{
	Name: "gato", Description: "Manda uma imagem aleatoria de uma gato",
	Aliases: []string{"cat", "meow"},
	Handler: func(ctx *base.CommandContext) {
		url, err := ImageService.GetRandomCat(ctx.Args != nil && strings.HasPrefix(ctx.Args[0], "gif"))

		if err != nil {
			ctx.Reply(Emotes.MIKU_CRY, sutils.Fmt("Um erro ocorreu ao procurar por um gato, desculpa. (`%s`)", err.Error()))
		} else {
			ctx.Send(url)
		}
	},
}
