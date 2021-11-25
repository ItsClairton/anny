package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var ShuffleCommand = core.Interaction{
	Name:        "embaralhar",
	Description: "Embaralhar as músicas da fila",
	Handler: func(ctx *core.InteractionContext) error {
		if ctx.VoiceState() == nil {
			return ctx.AsEphemeral().Send(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			return ctx.AsEphemeral().Send(emojis.Cry, "Não há nada tocando no momento.")
		}

		if len(player.Queue) < 2 {
			return ctx.AsEphemeral().Send(emojis.Cry, "Não há nada para embaralhar na fila.")
		}

		player.Shuffle()
		return ctx.Send(emojis.OK, "Músicas embaralhadas.")
	},
}
