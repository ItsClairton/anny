package anime

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/services/anilist"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/date"
	"github.com/ItsClairton/Anny/utils/i18n"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/ItsClairton/Anny/utils/sutils"
)

var MangaCommand = base.Command{
	Name: "manga",
	Handler: func(ctx *base.CommandContext) {

		if ctx.Args == nil {
			ctx.ReplyWithUsage("<nome de um mangá>")
			return
		}

		manga, err := anilist.SearchMediaAsManga(strings.Join(ctx.Args, " "))

		if err != nil {
			if err.Error() == "Not Found." {
				ctx.Reply(Emotes.MIKU_CRY, "anime.manga.not-found")
			} else {
				ctx.ReplyWithError(err)
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

		launchStr := sutils.Fmt("%s", date.ToPrettyDate(ctx.Locale, &manga.StartDate))

		if manga.Status == "NOT_YET_RELEASED" {
			launchStr = ctx.Locale.GetString("prevDate", launchStr)
		}

		if manga.EndDate.Year > 0 && manga.StartDate != manga.EndDate {
			launchStr = ctx.Locale.GetString("untilDate", launchStr, date.ToPrettyDate(ctx.Locale, &manga.EndDate))
		}

		hasTrailer := len(manga.GetTrailerURL()) > 0

		if hasTrailer {
			launchStr = sutils.Fmt("[%s](%s)", launchStr, manga.GetTrailerURL())
		}

		sourceStr := ctx.Locale.GetFromArray("anime.source", manga.GetSource())
		statusStr := ctx.Locale.GetFromArray("anime.status", manga.GetStatus())

		eb := embed.NewEmbed(ctx.Locale, "anime.manga.embed").
			SetTitle(sutils.Fmt("%s | %s", Emotes.HAPPY, manga.Title.JP)).
			SetDescription(rawSynopsis).
			SetURL(manga.SiteURL).
			SetThumbnail(manga.Cover.ExtraLarge).
			SetImage(manga.Banner).
			SetColor(sutils.ToHexNumber(manga.Cover.Color)).
			WithField(manga.GetCreator(), true).
			WithField(strings.Join(ctx.Locale.GetPrettyGenres(manga.Genres), ", "), true).
			WithField(strings.Join(manga.GetArts(), "\n"), true).
			WithField(chapters, true).
			WithField(sourceStr, true).
			WithField(volumes, true).
			WithField("N/A", true).
			WithField(launchStr, true).
			WithField(statusStr, true).
			SetFooter(sutils.Is(hasTrailer, ctx.Locale.GetString("anime.trailer-footer"), "Powered By AniList & MAL"), "https://anilist.co/img/icons/favicon-32x32.png")

		response := response.New(ctx.Locale).WithEmbed(eb)
		msg, err := ctx.ReplyWithResponse(response)

		if err != nil {
			logger.Warn(err.Error())
			return
		}

		if ctx.Locale.ID != "en_US" {
			translatedSynopsis, err := i18n.FromGoogle("auto", strings.Split(ctx.Locale.ID, "_")[0], rawSynopsis)

			if err == nil {
				eb.SetDescription(translatedSynopsis)
				ctx.EditWithResponse(msg.ID, response)
			}
		}

		mal, err := manga.GetBasicFromMAL()

		if err == nil {
			if mal.Score > 0 {
				eb.SetFieldValue(6, sutils.Fmt("%.2f", mal.Score))
			}

			if len(mal.Genres) > 0 {
				totalGenres := append(manga.Genres, mal.Genres...)
				eb.SetFieldValue(1, strings.Join(ctx.Locale.GetPrettyGenres(totalGenres), ", "))
			}
			ctx.EditWithResponse(msg.ID, response)
		}

	},
}
