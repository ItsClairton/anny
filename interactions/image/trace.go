package image

import (
	"bytes"
	"strings"

	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/providers"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/bwmarrin/discordgo"
)

var TraceContext = discord.Interaction{
	Name: "Que anime é esse?",
	Type: discordgo.MessageApplicationCommand,
	Handler: func(ctx *discord.InteractionContext) {
		message := ctx.ApplicationCommandData().Resolved.Messages[ctx.ApplicationCommandData().TargetID]
		attachment := getAttachment(message)

		if attachment == "" {
			ctx.ReplyWithEmote(emojis.MikuCry, "Não achei nenhuma imagem, GIF ou vídeo nessa mensagem.")
			return
		}
		response := discord.NewResponse().WithContentEmoji(emojis.AnimatedStaff, "Procurando...")
		err := ctx.SendResponse(response)
		if err != nil {
			logger.Warn("%s", err.Error())
			return
		}

		result, err := providers.SearchAnimeByScene(attachment)
		if err != nil {
			ctx.EditResponse(response.WithContentEmoji(emojis.MikuCry, "Um erro ocorreu ao executar esse comando. (`%s`)", err.Error()))
			return
		}

		ctx.EditResponse(response.WithContentEmoji(emojis.KannaPeer, "É uma cena%s de %s (%s).",
			utils.Is(result.Episode > 0, utils.Fmt(" do episódio **%d**", result.Episode), ""),
			utils.Is(len(result.Title.English) > 0 && !strings.EqualFold(result.Title.Japanese, result.Title.English),
				utils.Fmt("**%s** (**%s**)", result.Title.Japanese, result.Title.English),
				utils.Fmt("**%s**", result.Title.Japanese)),
			utils.Is(utils.ToDisplayTime(result.From) == utils.ToDisplayTime(result.To),
				utils.Fmt("`%s`", utils.ToDisplayTime(result.From)),
				utils.Fmt("`%s`/`%s`", utils.ToDisplayTime(result.From), utils.ToDisplayTime(result.To)))))

		video, err := utils.GetFromWeb(result.Video + "&size=l")

		if err == nil {
			response.WithContent("").WithFile(&discordgo.File{
				Name:        utils.Is(result.Adult, "SPOILER_preview.mp4", "preview.mp4"),
				ContentType: "video/mp4",
				Reader:      bytes.NewReader(video),
			})
			ctx.SendFollowUp(response)
		}
	},
}

func getAttachment(msg *discordgo.Message) string {
	if len(msg.Attachments) > 0 {
		return msg.Attachments[0].ProxyURL
	}

	if len(msg.Embeds) > 0 {
		if msg.Embeds[0].Image != nil {
			return msg.Embeds[0].Image.ProxyURL
		}
		if msg.Embeds[0].Thumbnail != nil {
			return msg.Embeds[0].Thumbnail.ProxyURL
		}
	}

	return ""
}
