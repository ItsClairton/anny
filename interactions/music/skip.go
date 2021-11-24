package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var SkipCommand = base.Interaction{
	Name:        "pular",
	Description: "Pular a música atual",
	Handler: func(ctx *base.InteractionContext) error {
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

		player.Skip()
		return ctx.Send(emojis.OK, "Música pulada com sucesso.")
	},
}
