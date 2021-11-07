package music

import (
	"time"

	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var NowplayingCommand = discord.Interaction{
	Name:        "tocando",
	Description: "Saber que música está tocando.",
	Handler: func(ctx *discord.InteractionContext) error {
		player := audio.GetPlayer(ctx.GuildID)

		if player == nil || player.Current == nil {
			return ctx.AsEphemeral().Send(emojis.MikuCry, "Não há nada tocando no momento.")
		}

		return ctx.WithEmbed(discord.NewEmbed().
			SetDescription("%s Tocando agora [%s](%s)", emojis.ZeroYeah, player.Current.Title, player.Current.URL).
			SetThumbnail(player.Current.Thumbnail).
			SetColor(0xA652BB).
			AddField("Autor", player.Current.Author, true).
			AddField("Duração", utils.Fmt("%v/%v", utils.FormatTime(player.Current.PlaybackPosition()), utils.Is(player.Current.IsLive, "--:--", utils.FormatTime(player.Current.Duration))), true).
			AddField("Provedor", player.Current.Provider(), true).
			SetFooter(utils.Fmt("Pedido por %s", player.Current.Requester.Username), player.Current.Requester.AvatarURL("")).
			SetTimestamp(player.Current.Time.Format(time.RFC3339))).Send()
	},
}
