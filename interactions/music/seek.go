package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

var SeekCommand = base.Interaction{
	Name:        "seek",
	Description: "Ir para um tempo especifico da música",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "posição",
		Description: "Posição desejada, formatos válidos são 01:22 ou 1m22s",
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

		if player.State == audio.PausedState {
			return ctx.AsEphemeral().Send(emojis.Cry, "Você precisa primeiro despausar a música para fazer isso.")
		}

		if player.Current.IsLive {
			return ctx.AsEphemeral().Send(emojis.Cry, "Você não pode fazer isso em transmissões ao vivo.")
		}

		duration, err := utils.ParseDuration(ctx.ArgumentAsString(0))
		if err != nil || duration < 0 || duration > player.Current.Duration {
			return ctx.AsEphemeral().Send(emojis.Cry, "Duração inválida ou maior do que a duração da música.")
		}

		defer player.Session.SetPosition(duration)
		return ctx.Send(emojis.Yeah, "Posição do player alterada para `%s`.", utils.FormatTime(duration))
	},
}
