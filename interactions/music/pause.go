package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PauseCommand = discord.Interaction{
	Name:        "pause",
	Description: "Pausar a música atual",
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
		if player.GetState() == audio.PausedState {
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "A música já está pausada.")
			return
		}
		player.Pause()
		ctx.ReplyWithEmote(emojis.PepeArt, "A música foi pausada com sucesso.")
	},
}
