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
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
			return
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.GetState() == audio.StoppedState {
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Não há nada tocando no momento.")
			return
		}
		queue := player.GetQueue()

		if len(queue) < 2 {
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Não há nada para embaralhar na fila.")
			return
		}

		player.Shuffle()
		ctx.ReplyWithEmote(emojis.ZeroYeah, "As músicas foram embaralhadas com sucesso.")
	},
}
