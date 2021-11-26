package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var NowplayingCommand = core.Interaction{
	Name:        "tocando",
	Description: "Saber que música está tocando.",
	Handler: func(ctx *core.InteractionContext) error {
		player := audio.GetPlayer(ctx.GuildID)

		if player == nil || player.Current == nil {
			return ctx.AsEphemeral().Send(emojis.Cry, "Não há nada tocando no momento.")
		}

		embed := core.NewEmbed().
			Description("%s Tocando no momento: **[%s](%s)**", emojis.AnimatedHype, player.Current.Title, player.Current.URL).
			Thumbnail(player.Current.Thumbnail).
			Color(0x00FF59).
			Field("Autor", player.Current.Author, true).
			Field("Duração", utils.Fmt("%v/%v", utils.FormatTime(player.Session.Position), utils.Is(player.Current.IsLive, "--:--", utils.FormatTime(player.Current.Duration))), true).
			Field("Provedor", player.Current.Provider(), true).
			Footer(utils.Fmt("Adicionado por %s#%s", player.Current.Requester.Username, player.Current.Requester.Discriminator), player.Current.Requester.AvatarURL()).
			Timestamp(player.Current.Time)

		if player.State == audio.PausedState {
			embed.Color(0xB4BE10).
				Description("%s Pausado no momento em: [%s](%s)", emojis.Cry, player.Current.Title, player.Current.URL)
		}

		return ctx.WithEmbed(embed).Send()
	},
}
