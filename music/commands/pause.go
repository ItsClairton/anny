package music

import (
	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PauseCommand = core.Command{
	Name:        "pausar",
	Description: "Pausar ou despausar a música atual",
	Handler:     func(ctx *core.CommandContext) { handleCommand(ctx) },
}

var ResumeCommand = core.Command{
	Name:        "despausar",
	Description: "Pausar ou despausar a música atual",
	Handler:     func(ctx *core.CommandContext) { handleCommand(ctx) },
}

func handleCommand(ctx *core.CommandContext) {
	if ctx.VoiceState() == nil {
		ctx.Ephemeral().Reply(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
		return
	}

	player := music.GetPlayer(ctx.GuildID)
	if player == nil || player.State == music.StoppedState {
		ctx.Ephemeral().Reply(emojis.Cry, "Não há nada tocando no momento.")
		return
	}

	if player.Current.IsLive {
		ctx.Ephemeral().Reply(emojis.Cry, "Você não pode fazer isso em transmissões ao vivo.")
		return
	}

	if player.State == music.PlayingState {
		player.Pause()
		ctx.Reply(emojis.OK, "Batidão pausada com sucesso.")
	} else {
		player.Resume()
		ctx.Reply(emojis.OK, "Batidão despausado com sucesso.")
	}
}
