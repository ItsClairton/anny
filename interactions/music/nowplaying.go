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
	Handler: func(ctx *discord.InteractionContext) {
		player := audio.GetPlayer(ctx.GuildID)
		if player == nil || player.State() == audio.StoppedState {
			ctx.SendEphemeral(emojis.MikuCry, "Não há nada tocando no momento.")
			return
		}

		current := player.Current()
		embed := discord.NewEmbed().
			SetDescription("%s Tocando agora [%s](%s)", emojis.ZeroYeah, current.Title, current.URL).
			SetThumbnail(current.Thumbnail).
			SetColor(0xA652BB).
			AddField("Autor", current.Author, true).
			AddField("Duração", utils.Fmt("%v/%v", utils.FormatTime(current.PlaybackPosition()), utils.Is(current.IsLive, "--:--", utils.FormatTime(current.Duration))), true).
			AddField("Provedor", current.Provider.Name(), true).
			SetFooter(utils.Fmt("Pedido por %s", current.Requester.Username), current.Requester.AvatarURL("")).
			SetTimestamp(current.Time.Format(time.RFC3339))

		ctx.SendEmbed(embed.Build())
	},
}
