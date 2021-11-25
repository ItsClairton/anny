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
			SetDescription("%s Tocando no momento: **[%s](%s)**", emojis.AnimatedHype, player.Current.Title, player.Current.URL).
			SetThumbnail(player.Current.Thumbnail).
			SetColor(0x00FF59).
			AddField("Autor", player.Current.Author, true).
			AddField("Duração", utils.Fmt("%v/%v", utils.FormatTime(player.Session.Position), utils.Is(player.Current.IsLive, "--:--", utils.FormatTime(player.Current.Duration))), true).
			AddField("Provedor", player.Current.Provider(), true).
			SetFooter(utils.Fmt("Adicionado por %s#%s", player.Current.Requester.Username, player.Current.Requester.Discriminator), player.Current.Requester.AvatarURL()).
			SetTimestamp(player.Current.Time)

		if player.State == audio.PausedState {
			embed.SetColor(0xB4BE10).
				SetDescription("%s Pausado no momento em: [%s](%s)", emojis.Cry, player.Current.Title, player.Current.URL)
		}

		return ctx.WithEmbed(embed).Send()
	},
}
