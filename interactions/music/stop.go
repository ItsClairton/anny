package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var StopCommand = core.Interaction{
	Name:        "parar",
	Description: "Parar a música atual, e limpar a fila",
	Handler: func(ctx *core.InteractionContext) error {
		if ctx.VoiceState() == nil {
			return ctx.AsEphemeral().Send(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.State == audio.StoppedState {
			return ctx.AsEphemeral().Send(emojis.Cry, "Não há nada tocando no momento.")
		}

		if player.State == audio.LoadingState {
			return ctx.AsEphemeral().Send(emojis.Cry, "Espere alguns segundos para fazer essa ação.")
		}

		player.Kill(true)
		return ctx.Send(emojis.Yeah, "Todas as músicas da fila foram limpas com sucesso.")
	},
}
