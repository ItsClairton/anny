package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var StopCommand = discord.Interaction{
	Name:        "parar",
	Description: "Parar a música atual, e limpar a fila",
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

		audio.RemovePlayer(player, true)
		ctx.ReplyWithEmote(emojis.ZeroYeah, "Música parada com sucesso, e fila limpa.")
	},
}
