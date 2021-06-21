package music

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/music"
)

var SkipCommand = base.Command{
	Name:    "pular",
	Aliases: []string{"skip", "s"},
	Handler: func(ctx *base.CommandContext) {

		player := music.GetPlayer(ctx.Message.GuildID)

		if player == nil || player.State != music.PlayingState {
			ctx.Reply(Emotes.ZERO_HMPF, "music.notPlaying")
			return
		}

		player.Current.Session.Source().Cleanup()
		ctx.Reply(Emotes.HAPPY, "music.skipSuccess")
	},
}
