package music

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/music"
)

var PauseCommand = base.Command{
	Name:    "pausar",
	Aliases: []string{"pause"},
	Handler: func(ctx *base.CommandContext) {

		player := music.GetPlayer(ctx.Message.GuildID)

		if player == nil || player.State == music.StoppedState {
			ctx.ReplyWithResponse(response.New(ctx.Locale).WithContentEmote(Emotes.ZERO_HMPF, "music.notPlaying"))
			return
		}

		isPaused := player.Current.Session.Paused()

		if isPaused {
			player.State = music.PlayingState
			ctx.ReplyWithResponse(response.New(ctx.Locale).WithContentEmote(Emotes.HAPPY, "music.unpausedSuccess"))
		} else {
			player.State = music.PausedState
			ctx.ReplyWithResponse(response.New(ctx.Locale).WithContentEmote(Emotes.HAPPY, "music.pausedSuccess"))
		}

		player.Current.Session.Pause()
	},
}
