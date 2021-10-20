package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PauseCommand = discord.Interaction{
	Name:        "pausar",
	Description: "Pausar a música atual",
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
		if player.State() == audio.PausedState {
			ctx.SendEphemeral(emojis.MikuCry, "A música já está pausada.")
			return
		}

		player.Pause()
		ctx.Send(emojis.PepeArt, "A música foi pausada com sucesso.")
	},
}
