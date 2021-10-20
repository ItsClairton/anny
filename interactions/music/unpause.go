package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var UnpauseCommand = discord.Interaction{
	Name:        "despausar",
	Description: "Despausar a música atual",
	Handler: func(ctx *discord.InteractionContext) {
		voiceId := ctx.GetVoiceChannel()
		if voiceId == "" {
			ctx.SendEphemeral(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
			return
		}
		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.State() == audio.StoppedState {
			ctx.SendEphemeral(emojis.MikuCry, "Não há nada tocando no momento.")
			return
		}
		if player.State() == audio.PlayingState {
			ctx.SendEphemeral(emojis.MikuCry, "A música já está despausada.")
			return
		}
		player.Unpause()
		ctx.Send(emojis.PepeArt, "A música foi despausada com sucesso.")
	},
}
