package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PauseCommand = discord.Interaction{
	Name:        "pausar",
	Description: "Pausar ou despausar a música atual",
	Handler:     handler,
}

var ResumeCommand = discord.Interaction{
	Name:        "despausar",
	Description: "Pausar ou despausar a música atual",
	Handler:     handler,
}

var handler = func(ctx *discord.InteractionContext) error {
	if ctx.VoiceState() == nil {
		return ctx.AsEphemeral().Send(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
	}

	player := audio.GetPlayer(ctx.GuildID)
	if player == nil || player.State() == audio.StoppedState {
		return ctx.AsEphemeral().Send(emojis.MikuCry, "Não há nada tocando no momento.")
	}

	if player.State() == audio.PausedState {
		player.Resume()
		return ctx.Send(emojis.PepeArt, "Música despausada com sucesso.")
	}

	player.Pause()
	return ctx.Send(emojis.PepeArt, "Música pausada com sucesso.")
}
