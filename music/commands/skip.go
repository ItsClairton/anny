package music

import (
	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var SkipCommand = core.Command{
	Name:        "pular",
	Description: "Pular a música atual",
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil || player.State == music.StoppedState {
			ctx.Ephemeral().Reply(emojis.Cry, "Não há nada para pular no momento.")
			return
		}

		player.Skip()
		ctx.Reply(emojis.OK, "Música pulada com sucesso.")
	},
}
