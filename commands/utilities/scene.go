package utilities

import (
	"bytes"
	"strings"
	"time"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/providers/image"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/constants"
	"github.com/bwmarrin/discordgo"
)

var SceneCommand = base.Command{
	Name:    "scene",
	Aliases: []string{"cena"},
	Handler: func(ctx *base.CommandContext) {

		var attachment string

		if ctx.MessageReference != nil {
			msg, err := ctx.Client.ChannelMessage(ctx.ChannelID, ctx.MessageReference.MessageID)
			if err == nil {
				attachment = utils.GetFirstAttachment(msg)
			}
		} else {
			attachment = utils.GetFirstAttachment(ctx.Message)
		}

		if attachment == "" {
			time.Sleep(300 * time.Millisecond)
			msg, err := ctx.Client.ChannelMessage(ctx.ChannelID, ctx.ID)

			if err == nil {
				attachment = utils.GetFirstAttachment(msg)
			}
		}

		if attachment != "" {

			response := response.New(ctx.Locale).WithContentEmote(constants.ANIMATED_STAFF, "searching")
			msg, _ := ctx.ReplyWithResponse(response)

			result, err := image.GetFromTrace(attachment)
			if err != nil {
				ctx.ReplyWithError(err)
			} else {
				var episodeStr string
				var titleStr string
				var timeStr string

				if len(result.Title.EN) > 0 && !strings.EqualFold(result.Title.JP, result.Title.EN) {
					titleStr = utils.Fmt("**%s** (**%s**)", result.Title.JP, result.Title.EN)
				} else {
					titleStr = utils.Fmt("**%s**", result.Title.JP)
				}

				if result.Episode > 0 {
					episodeStr = ctx.GetString("utilities.scene.ofEpisode", result.Episode, titleStr)
				} else {
					episodeStr = ctx.GetString("utilities.scene.of", titleStr)
				}

				fromTime := utils.ToHHMMSS(result.From)
				toTime := utils.ToHHMMSS(result.To)

				if fromTime != toTime {
					timeStr = ctx.GetString("utilities.scene.betweenMinutes", fromTime, toTime)
				} else {
					timeStr = ctx.GetString("utilities.scene.betweenMinute", fromTime)
				}

				finalResponse := ctx.GetString("utilities.scene.base", episodeStr, timeStr)
				response.SetContentEmote(constants.HAPPY, utils.Fmt("%s (%s)", finalResponse, ctx.GetString("utilities.scene.generatingPreview")))
				ctx.EditWithResponse(msg.ID, response)

				videoBody, err := utils.GetFromWeb(result.Video + "&size=l")

				if err != nil {
					response.SetContentEmote(constants.MIKU_CRY, utils.Fmt("%s (%s)", finalResponse, ctx.GetString("utilities.scene.previewError")))
					ctx.EditWithResponse(msg.ID, response)
					return
				}

				ctx.ReplyWithResponse(response.SetContentEmote(constants.YEAH, finalResponse).WithFile(&discordgo.File{
					Name:        utils.Is(result.Adult, "SPOILER_preview.mp4", "preview.mp4"),
					ContentType: "mp4",
					Reader:      bytes.NewReader(videoBody),
				}))

				ctx.DeleteMessage(msg)
			}

		} else {
			ctx.Reply(constants.ZERO_HMPF, "utilities.scene.usage")
		}

	},
}
