package music

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/music"
	"github.com/bwmarrin/discordgo"
)

func getVoice(userId string, guild *discordgo.Guild) string {

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userId {
			return vs.ChannelID
		}
	}

	return ""
}

var PlayCommand = base.Command{
	Name:    "tocar",
	Aliases: []string{"play", "p"},
	Handler: func(ctx *base.CommandContext) {
		if ctx.Args == nil {
			ctx.ReplyWithUsage("<nome de uma música>")
			return
		}
		g, _ := ctx.Client.State.Guild(ctx.Message.GuildID)
		channelId := getVoice(ctx.Author.ID, g)
		if len(channelId) > 0 {
			var msg *discordgo.Message
			var track *music.Track
			var err error
			content := strings.Join(ctx.Args, " ")
			response := response.New(ctx.Locale)
			if music.YouTubeRegex.MatchString(content) {
				response.WithContentEmote(Emotes.ANIMATED_STAFF, "music.detail")
				msg, _ = ctx.ReplyWithResponse(response)
				track, err = music.GetTrackFromYouTube(content)
			} else {
				response.WithContentEmote(Emotes.ANIMATED_STAFF, "searching")
				msg, err = ctx.ReplyWithResponse(response)
				id, err := music.GetIDFromYouTube(content)
				if err == nil {
					if id == "" {
						ctx.EditWithResponse(msg.ID, response.WithContentEmote(Emotes.ZERO_HMPF, "music.notFound"))
						return
					} else {
						ctx.EditWithResponse(msg.ID, response.WithContentEmote(Emotes.ANIMATED_STAFF, "music.detail"))
						track, err = music.GetTrackFromYouTube(id)
						if err != nil {
							ctx.EditWithResponse(msg.ID, response.WithContentEmote(Emotes.MIKU_CRY, "error", err.Error()))
							return
						}
					}
				}
			}
			if err != nil {
				ctx.EditWithResponse(msg.ID, response.WithContentEmote(Emotes.MIKU_CRY, "error", err.Error()))
				return
			}

			if track == nil {
				return
			}
			track.Requester = ctx.Author
			player := music.GetPlayer(g.ID)
			if player == nil {
				vc, err := ctx.Client.ChannelVoiceJoin(g.ID, channelId, false, true)
				if err != nil {
					ctx.ReplyWithError(err)
					return
				}
				player = music.AddPlayer(&music.Player{State: music.StoppedState, Guild: g, Connection: vc, Ctx: ctx})
			}
			response.WithContentEmote(Emotes.HAPPY, "music.addedQueue", track.Name, track.Author, len(player.Tracks)+1)
			ctx.EditWithResponse(msg.ID, response)
			player.LoadTrack(*track)
		} else {
			ctx.Reply(Emotes.ZERO_HMPF, "music.notConnected")
		}
	},
}
