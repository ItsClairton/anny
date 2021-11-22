package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

var VolumeCommand = base.Interaction{
	Name:        "volume",
	Description: "Alterar o volume da música",
	Options: discord.CommandOptions{&discord.IntegerOption{
		OptionName:  "volume",
		Description: "Volume da música em porcentagem",
		Required:    true,
	}},
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

		argument := ctx.ArgumentAsInteger(0)
		if argument <= 0 || argument > 100 {
			return ctx.AsEphemeral().Send(emojis.Cry, "Volume invalido, você só pode alterar o volume entre 1 a 100%%.")
		}

		player.Session.SetVolume(argument)
		return ctx.Send(emojis.Yeah, "Volume alterado para %d%% com sucesso.", argument)
	},
}
