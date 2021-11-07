package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var StopCommand = discord.Interaction{
	Name:        "parar",
	Description: "Parar a música atual, e limpar a fila",
	Handler: func(ctx *discord.InteractionContext) error {
		if ctx.VoiceState() == nil {
			return ctx.AsEphemeral().Send(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.State == audio.StoppedState {
			return ctx.AsEphemeral().Send(emojis.MikuCry, "Não há nada tocando no momento.")
		}

		player.Kill(true)
		return ctx.Send(emojis.ZeroYeah, "Todas as músicas da fila foram limpas com sucesso.")
	},
}
