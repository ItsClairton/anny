package music

import (
	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var ShuffleCommand = core.Command{
	Name:        "embaralhar",
	Description: "Embaralhar as músicas da fila",
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil || player.State == music.StoppedState {
			ctx.Ephemeral().Reply(emojis.Cry, "Não há nada tocando no momento.")
			return
		}

		if len(player.Queue) < 2 {
			ctx.Ephemeral().Reply(emojis.Cry, "Não há músicas suficientes para embaralhar na fila.")
			return
		}

		player.Shuffle()
		ctx.Reply(emojis.OK, "Músicas embaralhadas com sucesso.")
	},
}
