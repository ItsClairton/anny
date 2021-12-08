package music

import (
	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

var SeekCommand = core.Command{
	Name:        "seek",
	Description: "Alterar a posição do batidão",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "posição",
		Description: "Posição desejada, exemplo de formatos válidos: 05:05 ou 5m5s",
		Required:    true,
	}},
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil || player.State != music.PlayingState {
			ctx.Ephemeral().Reply(emojis.Cry, "Não há nada tocando no momento ou o batidão está pausado.")
			return
		}

		if player.Current.IsLive {
			ctx.Ephemeral().Reply(emojis.Cry, "Você não pode fazer isso em transmissões ao vivo.")
			return
		}

		duration, err := utils.ParseDuration(ctx.Argument(0).String())
		if err != nil || duration < 0 || duration > player.Current.Duration {
			ctx.Ephemeral().Reply(emojis.Cry, "Duração inválida ou maior que a duração total da música.")
			return
		}

		player.Voicy.Seek(duration)
		ctx.Reply(emojis.OK, "Posição do batidão alterada para os minutos `%s`.", utils.FormatTime(duration))
	},
}
