package utilities

import (
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/i18n"
	"github.com/ItsClairton/Anny/logger"
	"github.com/ItsClairton/Anny/services/anilist"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/constants"
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
			ctx.ReplyWithError(err)
			return
		}

		if manga == nil {
			ctx.Reply(constants.MIKU_CRY, "utilities.manga.notFound")
			return
		}

		rawSynopsis := utils.ToMD(manga.Synopsis)
		chapters := utils.Fmt("%d", manga.Chapters)
		volumes := utils.Fmt("%d", manga.Volumes)

		if manga.Chapters < 1 {
			chapters = "N/A"
		}

		if manga.Volumes < 1 {
			volumes = "N/A"
		}

		launchStr := utils.Fmt("%s", ctx.ToPrettyDate(manga.StartDate))

		if manga.Status == "NOT_YET_RELEASED" {
			launchStr = ctx.GetString("prevDate", launchStr)
		}

		if manga.EndDate.Year > 0 && manga.StartDate != manga.EndDate {
			launchStr = ctx.GetString("untilDate", launchStr, ctx.ToPrettyDate(manga.EndDate))
		}

		hasTrailer := len(manga.GetTrailerURL()) > 0

		if hasTrailer {
			launchStr = utils.Fmt("[%s](%s)", launchStr, manga.GetTrailerURL())
		}

		sourceStr := ctx.GetFromArray("utilities.source", manga.GetSource())
		statusStr := ctx.GetFromArray("utilities.status", manga.GetStatus())

		eb := embed.NewEmbed(ctx.Locale, "utilities.manga.embed").
			SetTitle(utils.Fmt("%s | %s", constants.HAPPY, manga.Title.JP)).
			SetDescription(rawSynopsis).
			SetURL(manga.SiteURL).
			SetThumbnail(manga.Cover.ExtraLarge).
			SetImage(manga.Banner).
			SetColor(utils.ToHexNumber(manga.Cover.Color)).
			WithField(manga.GetCreator(), true).
			WithField(strings.Join(ctx.GetPrettyGenres(manga.Genres), ", "), true).
			WithField(strings.Join(manga.GetArts(), "\n"), true).
			WithField(chapters, true).
			WithField(sourceStr, true).
			WithField(volumes, true).
			WithField("N/A", true).
			WithField(launchStr, true).
			WithField(statusStr, true).
			SetFooter(utils.Is(hasTrailer, ctx.GetString("utilities.trailer-footer"), "Powered By AniList & MAL"), "https://anilist.co/img/icons/favicon-32x32.png")

		msg, err := ctx.ReplyWithEmbed(eb)

		if err != nil {
			logger.Warn(err.Error())
			return
		}

		if ctx.ID != "en_US" {
			translatedSynopsis, err := i18n.FromGoogle("auto", strings.Split(ctx.ID, "_")[0], rawSynopsis)

			if err == nil {
				eb.SetDescription(translatedSynopsis)
				ctx.EditWithEmbed(msg.ID, eb)
			}
		}

		mal, err := manga.GetBasicFromMAL()

		if err == nil {
			if mal.Score > 0 {
				eb.SetFieldValue(6, utils.Fmt("%.2f", mal.Score))
			}

			if len(mal.Genres) > 0 {
				totalGenres := append(manga.Genres, mal.Genres...)
				eb.SetFieldValue(1, strings.Join(ctx.GetPrettyGenres(totalGenres), ", "))
			}
			ctx.EditWithEmbed(msg.ID, eb)
		}

	},
}
