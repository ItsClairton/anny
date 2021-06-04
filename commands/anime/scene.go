package anime

import (
	"bytes"
	"time"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/services/image"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/rest"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/bwmarrin/discordgo"
)

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
	msg, _ := ctx.Reply(Emotes.ANIMATED_STAFF, "Obtendo resultados...")

	result, err := image.GetFromTrace(attachment)
	if err != nil {
		ctx.EditReply(msg, Emotes.MIKU_CRY, sutils.Fmt("Um erro ocorreu ao entrar em contato com o trace.moe, dsclpa. (`%s`)", err))
	} else {
		response := "Talvez seja uma cena"

		if result.Episode > 0 {
			response += sutils.Fmt(" do episódio **%d** de", result.Episode)
		} else {
			response += sutils.Fmt(" de")
		}

		if result.Title.EN != nil && sutils.ToLower(result.Title.JP) != sutils.ToLower(result.Title.EN) {
			response += sutils.Fmt(" **%s** (**%s**)", result.Title.JP, result.Title.EN)
		} else {
			response += sutils.Fmt(" **%s**", result.Title.JP)
		}

		fromTime := sutils.ToHHMMSS(result.From)
		toTime := sutils.ToHHMMSS(result.To)

		if fromTime != toTime {
			response += sutils.Fmt(" que aparece entre os minutos `%s` e `%s`", fromTime, toTime)
		} else {
			response += sutils.Fmt(" que aparece no minuto `%s`", fromTime)
		}

		ctx.EditReply(msg, Emotes.YEAH, sutils.Fmt("%s. (Gerando Preview)", response))

		videoBody, err := rest.Get(result.Video + "&size=l")

		if err != nil {
			ctx.EditReply(msg, Emotes.YEAH, sutils.Fmt("%s. (Um erro ocorreu ao tentar gerar um Preview)", response))
			return
		}

		ctx.ReplyWithFile(Emotes.YEAH, sutils.Fmt("%s.", response), &discordgo.File{
			Name:        "preview.mp4",
			ContentType: "mp4",
			Reader:      bytes.NewReader(videoBody),
		})
		ctx.DeleteMessage(msg)
	}
}

var SceneCommand = base.Command{
	Name: "cena", Description: "Saber qual anime, em qual episódio e quais minutos uma cena especifica (Foto, GIF ou vídeo) apareceu",
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
				ctx.Reply(Emotes.MIKU_CRY, "Você precisa mandar um link ou anexar, uma imagem, um GIF, ou um vídeo. (Lembrando que GIF e Vídeo só é analisado o primeiro frame), ou então você pode referenciar uma mensagem que contenha algum destes itens.")
			}

		} else {
			sendTraceMessage(ctx, attachment)
		}
	},
}
