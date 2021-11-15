package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var ShuffleCommand = base.Interaction{
	Name:        "embaralhar",
	Description: "Embaralhar as músicas da fila",
	Handler: func(ctx *base.InteractionContext) error {
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
		return ctx.Send(emojis.Yeah, "As músicas foram embaralhadas com sucesso.")
	},
}
