package anime

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/services/anilist"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/embed"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/ItsClairton/Anny/utils/translate"
)

var MangaCommand = base.Command{
	Name: "manga", Description: "Saber informações básicas sobre um mangá",
	Handler: func(ctx *base.CommandContext) {

		if ctx.Args == nil {
			ctx.Reply(Emotes.MIKU_CRY, "Você precisa falar o nome do mangá")
			return
		}

		manga, err := anilist.SearchMediaAsManga(strings.Join(ctx.Args, " "))

		if err != nil {
			if err.Error() == "Not Found." {
				ctx.Reply(Emotes.MIKU_CRY, "Não encontrei informações sobre esse mangá, Desculpa ;(")
			} else {
				ctx.Reply(Emotes.MIKU_CRY, sutils.Fmt("Houve um erro ao obter informações sobre esse mangá, desculpa. (%s)", err.Error()))
			}
			return
		}

		rawSynopsis := sutils.ToMD(manga.Synopsis)
		chapters := sutils.Fmt("%d", manga.Chapters)
		volumes := sutils.Fmt("%d", manga.Volumes)

		if manga.Chapters < 1 {
			chapters = "N/A"
		}

		if manga.Volumes < 1 {
			volumes = "N/A"
		}

		launchStr := sutils.Fmt("%s", manga.GetPrettyStartDate())

		if manga.EndDate.Year > 0 && manga.StartDate != manga.EndDate {
			launchStr += sutils.Fmt("\n%s", manga.GetPrettyEndDate())
		}

		if len(manga.GetTrailerURL()) > 0 {
			launchStr = sutils.Fmt("[%s](%s)", launchStr, manga.GetTrailerURL())
		}

		eb := embed.NewEmbed().
			SetTitle(sutils.Fmt("📰 %s", manga.Title.JP)).
			SetDescription(rawSynopsis).
			SetURL(manga.SiteURL).
			SetThumbnail(manga.Cover.ExtraLarge).
			SetImage(manga.Banner).
			AddField("História", manga.GetCreator(), true).
			AddField("Status", manga.GetPrettyStatus(), true).
			AddField("Arte", strings.Join(manga.GetArts(), "\n"), true).
			AddField("Capitulos", chapters, true).
			AddField("Gêneros", strings.Join(manga.Genres, ", "), true).
			AddField("Volumes", volumes, true).
			AddField("Pontuação", "...", true).
			AddField("Data de Lançamento", launchStr, true).
			AddField("Adaptação", manga.GetPrettySource(), true).
			SetFooter("Powered By AniList & MAL", "https://anilist.co/img/icons/favicon-32x32.png")

		color, err := sutils.ToHexNumber(manga.Cover.Color)

		if err == nil {
			eb.SetColor(color)
		}

		msg, err := ctx.ReplyWithEmbed(eb.MessageEmbed)

		if err != nil {
			logger.Warn(err.Error())
			return
		}

		translatedSynopsis, err := translate.Translate("auto", "pt", rawSynopsis)

		if err == nil {
			eb.SetDescription(translatedSynopsis)
		}

		translatedGenres, err := translate.Translate("auto", "pt", strings.Join(manga.Genres, ", "))

		if err == nil {
			array := strings.Split(translatedGenres, ", ")
			var newArray []string
			for _, t := range array {
				if strings.Contains(strings.ToLower(t), "fatia") {
					newArray = append(newArray, "Slice of Life")
				} else {
					newArray = append(newArray, strings.Title(t))
				}
			}

			eb.SetField(4, "Gêneros", strings.Join(newArray, ", "), true)
		}

		ctx.EditWithEmbed(msg, eb.MessageEmbed)

		score, err := manga.GetScoreFromMAL()

		if err == nil {
			eb.SetField(6, "Pontuação", sutils.Fmt("%.2f", score), true)
			ctx.EditWithEmbed(msg, eb.MessageEmbed)
		}

	},
}
