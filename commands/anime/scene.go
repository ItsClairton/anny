package anime

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/services/image"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/constants"
	"github.com/bwmarrin/discordgo"
)

var SceneCommand = base.Command{
	Name: "cena",
	Handler: func(ctx *base.CommandContext) {

		attachment := ""

		if ctx.Message.MessageReference != nil {
			ref, err := ctx.Client.ChannelMessage(ctx.Message.ChannelID, ctx.Message.MessageReference.MessageID)
			if err == nil {
				attachment = getURLFromMessage(ref)
			}
		}

		if len(attachment) < 1 {
			time.Sleep(300 * time.Millisecond)
			refresh, _ := ctx.Client.ChannelMessage(ctx.Message.ChannelID, ctx.Message.ID) // Tem de pegar a mensagem de novo porque o Discord demora renderizar o embed as vezes.
			attachment = getURLFromMessage(refresh)

			if len(attachment) > 1 {
				sendTraceMessage(ctx, attachment)
			} else {
				ctx.Reply(constants.MIKU_CRY, "anime.scene.usage")
			}
		} else {
			sendTraceMessage(ctx, attachment)
		}
	},
}

func getURLFromMessage(msg *discordgo.Message) string {

	if len(msg.Attachments) > 0 {
		return msg.Attachments[0].ProxyURL
	}
	if len(msg.Embeds) > 0 && (msg.Embeds[0].Image != nil || msg.Embeds[0].Thumbnail != nil) {
		if msg.Embeds[0].Image != nil {
			return msg.Embeds[0].Image.ProxyURL
		}
		return msg.Embeds[0].Thumbnail.ProxyURL
	}

	return ""
}

func sendTraceMessage(ctx *base.CommandContext, attachment string) {
	response := response.New(ctx.Locale).WithContentEmote(constants.ANIMATED_STAFF, "searching")
	msg, _ := ctx.ReplyWithResponse(response)

	result, err := image.GetFromTrace(attachment)
	if err != "" {
		ctx.ReplyWithError(errors.New(utils.Fmt("trace error: %s", err)))
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
			episodeStr = ctx.GetString("anime.scene.ofEpisode", result.Episode, titleStr)
		} else {
			episodeStr = ctx.GetString("anime.scene.of", titleStr)
		}

		fromTime := utils.ToHHMMSS(result.From)
		toTime := utils.ToHHMMSS(result.To)

		if fromTime != toTime {
			timeStr = ctx.GetString("anime.scene.betweenMinutes", fromTime, toTime)
		} else {
			timeStr = ctx.GetString("anime.scene.betweenMinute", fromTime)
		}

		finalResponse := ctx.GetString("anime.scene.base", episodeStr, timeStr)
		response.SetContentEmote(constants.HAPPY, utils.Fmt("%s (%s)", finalResponse, ctx.GetString("anime.scene.generatingPreview")))
		ctx.EditWithResponse(msg.ID, response)

		videoBody, err := utils.GetFromWeb(result.Video + "&size=l")

		if err != nil {
			response.SetContentEmote(constants.MIKU_CRY, utils.Fmt("%s (%s)", finalResponse, ctx.GetString("anime.scene.previewError")))
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
}
