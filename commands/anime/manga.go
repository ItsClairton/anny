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
			launchStr += sutils.Fmt(" até %s", manga.GetPrettyEndDate())
		}

		hasTrailer := len(manga.GetTrailerURL()) > 0

		if hasTrailer {
			launchStr = sutils.Fmt("[%s](%s)", launchStr, manga.GetTrailerURL())
		}

		eb := embed.NewEmbed().
			SetAuthor(sutils.Fmt("Tipo: %s - Episódios: %d", manga.GetPrettyFormat(), manga.Episodes), "https://cdn.discordapp.com/avatars/743538534589267990/a6c5e905673d041a88b49203d6bc74dd.png?size=2048").
			SetTitle(sutils.Fmt("%s | %s", Emotes.HAPPY, manga.Title.JP)).
			SetDescription(rawSynopsis).
			SetURL(manga.SiteURL).
			SetThumbnail(manga.Cover.ExtraLarge).
			SetImage(manga.Banner).
			SetColor(sutils.ToHexNumber(manga.Cover.Color)).
			AddField("História", manga.GetCreator(), true).
			AddField("Gênero", strings.Join(manga.GetPrettyGenres(), ", "), true).
			AddField("Ilustração", strings.Join(manga.GetArts(), "\n"), true).
			AddField("Capitulos", chapters, true).
			AddField("Adaptação", manga.GetPrettySource(), true).
			AddField("Volumes", volumes, true).
			AddField("Pontuação", "N/A", true).
			AddField("Data de Estreia", launchStr, true).
			AddField("Status", manga.GetPrettyStatus(), true).
			SetFooter(sutils.Is(hasTrailer, "Clique na data de estreia para ver o Trailer", "Powered By AniList & MAL"), "https://anilist.co/img/icons/favicon-32x32.png")

		msg, err := ctx.ReplyWithEmbed(eb.Build())

		if err != nil {
			logger.Warn(err.Error())
			return
		}

		translatedSynopsis, err := translate.Translate("auto", "pt", rawSynopsis)

		if err == nil {
			eb.SetDescription(translatedSynopsis)
		}

		ctx.EditWithEmbed(msg, eb.Build())

		mal, err := manga.GetBasicFromMAL()

		if err == nil {
			if mal.Score > 0 {
				eb.SetField(6, "Pontuação", sutils.Fmt("%.2f", mal.Score), true)
			}

			if len(mal.Genres) > 0 {
				translatedGenres, err := translate.Translate("en", "pt", strings.Join(mal.Genres, ", "))

				if err == nil {
					var finalGenres []string
					finalGenres = manga.GetPrettyGenres()

					for _, genre := range strings.Split(translatedGenres, ", ") {
						finalGenres = append(finalGenres, strings.Title(genre))
					}

					eb.SetField(1, "Gênero", strings.Join(finalGenres, ", "), true)
				}
			}
			ctx.EditWithEmbed(msg, eb.Build())
		}

	},
}
