package music

import (
	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var NowplayingCommand = core.Command{
	Name:        "tocando",
	Description: "Mostra que música que está tocando agora",
	Handler: func(ctx *core.CommandContext) {
		player := music.GetPlayer(ctx.GuildID)

		if player.State == music.StoppedState {
			ctx.Ephemeral().Reply(emojis.Cry, "Você não pode fazer isso em transmissões ao vivo.")
			return
		}

		embed := utils.NewEmbed().
			Description("%s Tocando no momento: **[%s](%s)**", emojis.AnimatedHype, player.Current.Title, player.Current.URL).
			Thumbnail(player.Current.Thumbnail).
			Color(0x00FF59).
			Field("Autor", player.Current.Author, true).
			Field("Duração", utils.Fmt("%v/%v", utils.FormatTime(player.Voicy.Position), utils.Is(player.Current.IsLive, "--:--", utils.FormatTime(player.Current.Duration))), true).
			Field("Provedor", player.Current.Provider(), true).
			Footer(utils.Fmt("Adicionado por %s#%s", player.Current.Requester.Username, player.Current.Requester.Discriminator), player.Current.Requester.AvatarURL()).
			Timestamp(player.Current.RequestedAt)

		if player.State == music.PausedState {
			embed.Color(0xB4BE10).Description("%s Pausado no momento em: [%s](%s)", emojis.Cry, player.Current.Title, player.Current.URL)
		}

		ctx.Reply(embed)
	},
}
