package image

import (
	"bytes"
	"strings"

	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/providers"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

var TraceContext = discord.Interaction{
	Name:     "Que anime é esse?",
	Type:     discordgo.MessageApplicationCommand,
	Deffered: true,
	Handler: func(ctx *discord.InteractionContext) error {
		message := ctx.ApplicationCommandData().Resolved.Messages[ctx.ApplicationCommandData().TargetID]

		attachment := getAttachment(message)
		if attachment == "" {
			return ctx.Send(emojis.MikuCry, "Não achei nenhuma imagem, GIF ou vídeo nessa mensagem.")
		}

		result, err := providers.SearchAnimeByScene(attachment)
		if err != nil {
			return ctx.SendWithError(err)
		}

		content := utils.Fmt("Talvez seja uma cena (%s)%s de %s.",
			utils.Is(result.From == result.To,
				utils.Fmt("`%s`", utils.FormatTime(result.From)),
				utils.Fmt("`%s`/`%s`", utils.FormatTime(result.From), utils.FormatTime(result.To))),
			utils.Is(result.Episode > 0, utils.Fmt(" do episódio **%d**", result.Episode), ""),
			utils.Is(len(result.Title.English) > 0 && !strings.EqualFold(result.Title.Japanese, result.Title.English),
				utils.Fmt("**%s** (**%s**)", result.Title.Japanese, result.Title.English),
				utils.Fmt("**%s**", result.Title.Japanese)))
		ctx.Send(emojis.KannaPeer, "%s (Gerando Preview)", content)

		video, err := utils.GetFromWeb(result.Video + "&size=l")
		if err != nil {
			return ctx.Edit(emojis.KannaPeer, "%s (Não foi possivel gerar o Preview)", content)
		}

		return ctx.WithFile(&discordgo.File{
			Name:        utils.Is(result.Adult, "SPOILER_preview.mp4", "preview.mp4"),
			ContentType: "video/mp4",
			Reader:      bytes.NewReader(video),
		}).Edit(emojis.KannaPeer, content)
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
