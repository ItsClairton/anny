package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var ShuffleCommand = discord.Interaction{
	Name:        "embaralhar",
	Description: "Embaralhar as músicas da fila",
	Handler: func(ctx *discord.InteractionContext) {
		if ctx.GetVoiceChannel() == "" {
			ctx.SendEphemeral(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
			return
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.State() == audio.StoppedState {
			ctx.SendEphemeral(emojis.MikuCry, "Não há nada tocando no momento.")
			return
		}
		queue := player.Queue()

		if len(queue) < 2 {
			ctx.SendEphemeral(emojis.MikuCry, "Não há nada para embaralhar na fila.")
			return
		}

		player.Shuffle()
		ctx.Send(emojis.ZeroYeah, "As músicas foram embaralhadas com sucesso.")
	},
}
