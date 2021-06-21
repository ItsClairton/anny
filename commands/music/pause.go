package music

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/music"
)

var PauseCommand = base.Command{
	Name:    "pausar",
	Aliases: []string{"pause"},
	Handler: func(ctx *base.CommandContext) {

		player := music.GetPlayer(ctx.Message.GuildID)

		if player == nil || player.State == music.StoppedState {
			ctx.Reply(Emotes.ZERO_HMPF, "music.notPlaying")
			return
		}

		isPaused := player.Current.Session.Paused()

		if isPaused {
			player.State = music.PlayingState
			ctx.Reply(Emotes.HAPPY, "music.unpausedSuccess")
		} else {
			player.State = music.PausedState
			ctx.Reply(Emotes.HAPPY, "music.pausedSuccess")
		}

		player.Current.Session.Pause()
	},
}
