package music

import (
	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var StopCommand = core.Command{
	Name:        "parar",
	Description: "Parar a música atual, limpar a fila e desconectar do canal de voz.",
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil {
			ctx.Ephemeral().Reply(emojis.Cry, "Não há nada tocando no momento.")
			return
		}

		player.Stop(false)
		ctx.Reply(emojis.OK, "Batidão parado com sucesso.")
	},
}
