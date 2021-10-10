package music

import (
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
		if player == nil || player.GetState() == audio.StoppedState {
			ctx.SendEphemeral(emojis.MikuCry, "Não há nada tocando no momento.")
			return
		}

		current := player.GetCurrent()
		ctx.SendEmbed(discord.NewEmbed().
			SetColor(0x0099E1).
			SetDescription(utils.Fmt("[%s](%s)", current.Title, current.URL)).
			SetThumbnail(current.ThumbnailUrl).
			AddField("Autor", player.GetCurrent().Author, true).
			AddField("Duração", utils.Fmt("%s/%s",
				utils.ToDisplayTime(current.Session.PlaybackPosition().Seconds()),
				utils.ToDisplayTime(current.Duration.Seconds())), true).Build())
	},
}
