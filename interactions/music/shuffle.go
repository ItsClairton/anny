package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var ShuffleCommand = discord.Interaction{
	Name:        "embaralhar",
	Description: "Embaralhar as músicas da fila",
	Handler: func(ctx *discord.InteractionContext) error {
		if ctx.VoiceState() == nil {
			return ctx.AsEphemeral().Send(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.State == audio.StoppedState {
			return ctx.AsEphemeral().Send(emojis.MikuCry, "Não há nada tocando no momento.")
		}

		if len(player.Queue) < 2 {
			return ctx.AsEphemeral().Send(emojis.MikuCry, "Não há nada para embaralhar na fila.")
		}

		player.Shuffle()
		return ctx.Send(emojis.ZeroYeah, "As músicas foram embaralhadas com sucesso.")
	},
}
