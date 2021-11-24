package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PauseCommand = base.Interaction{
	Name:        "pausar",
	Description: "Pausar ou despausar a música atual",
	Handler:     handler,
}

var ResumeCommand = base.Interaction{
	Name:        "despausar",
	Description: "Pausar ou despausar a música atual",
	Handler:     handler,
}

var handler = func(ctx *base.InteractionContext) error {
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

	if player.Current.IsLive {
		return ctx.AsEphemeral().Send(emojis.Cry, "Você não pode pausar transmissões ao vivo.")
	}

	if player.State == audio.PausedState {
		player.Resume()
		return ctx.Send(emojis.OK, "Música despausada.")
	}

	player.Pause()
	return ctx.Send(emojis.OK, "Música pausada.")
}
