package music

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/constants"
	"github.com/ItsClairton/Anny/utils/music"
	"github.com/ItsClairton/Anny/utils/music/provider"
)

var PlayCommand = base.Command{
	Name:    "tocar",
	Aliases: []string{"play", "p"},
	Handler: func(ctx *base.CommandContext) {

		if ctx.Args == nil {
			ctx.ReplyWithUsage("<nome de uma música>")
			return
		}

		voiceId := ctx.GetVoice()
		if voiceId != "" {

			msg, err := ctx.Reply(constants.KANNAPEER, "searching")
			if err != nil {
				return
			}

			info, err := provider.GetInfo(strings.Join(ctx.Args, " "))
			if err != nil {
				ctx.ReplyWithError(err)
				return
			}

			if info == nil {
				ctx.Edit(msg.ID, constants.MIKU_CRY, "music.notFound")
				return
			}

			go ctx.Edit(msg.ID, constants.TOHRU, "music.addedQueue", info.Title, info.Author)

			player := music.GetPlayer(ctx.Message.GuildID)
			if player == nil {
				vc, err := ctx.Client.ChannelVoiceJoin(ctx.Message.GuildID, voiceId, false, true)
				if err != nil {
					ctx.ReplyWithError(err)
					return
				}

				player = music.AddPlayer(&music.Player{
					State:      music.StoppedState,
					GuildID:    ctx.Message.GuildID,
					Connection: vc,
					Ctx:        ctx,
				})
			}

			player.AddQueue(music.Track{
				PartialInfo: info,
				Stream:      nil,
				Requester:   ctx.Author,
			})
		} else {
			ctx.Reply(constants.MIKU_CRY, "music.notConnected")
		}

	},
}
