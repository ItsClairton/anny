package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var SkipCommand = discord.Interaction{
	Name:        "pular",
	Description: "Pular a música atual que está tocando",
	Handler: func(ctx *discord.InteractionContext) {
		voiceId := ctx.GetVoiceChannel()
		if voiceId == "" {
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
			return
		}
		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.GetState() == audio.StoppedState {
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Não há nada tocando no momento.")
			return
		}

		player.Skip()
		ctx.ReplyWithEmote(emojis.PepeArt, "Música pulada com sucesso.")
	},
}
