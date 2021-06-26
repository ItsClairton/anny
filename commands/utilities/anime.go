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

var AnimeCommand = base.Command{
	Name: "anime",
	Handler: func(ctx *base.CommandContext) {

		if ctx.Args == nil {
			ctx.ReplyWithUsage("<nome de um anime>")
			return
		}

		anime, err := anilist.SearchMediaAsAnime(strings.Join(ctx.Args, " "))

		if err != nil {
			ctx.ReplyWithError(err)
			return
		}

		if anime == nil {
			ctx.Reply(constants.MIKU_CRY, "utilities.anime.notFound")
			return
		}

		launchStr := utils.Fmt("%s", ctx.ToPrettyDate(anime.StartDate))

		if anime.Status == "NOT_YET_RELEASED" {
			launchStr = ctx.GetString("prevDate", launchStr)
		}

		if anime.EndDate.Year > 0 && anime.StartDate != anime.EndDate {
			launchStr = ctx.GetString("untilDate", launchStr, ctx.ToPrettyDate(anime.EndDate))
		}

		hasTrailer := len(anime.GetTrailerURL()) > 0

		if hasTrailer {
			launchStr = utils.Fmt("[%s](%s)", launchStr, anime.GetTrailerURL())
		}

		rawSynopsis := utils.ToMD(anime.Synopsis)

		if err != nil {
			ctx.ReplyWithError(err)
			return
		}

		typeStr := ctx.GetFromArray("utilities.type", anime.GetType())
		sourceStr := ctx.GetFromArray("utilities.source", anime.GetSource())
		seasonStr := ctx.GetFromArray("utilities.season", anime.GetSeason())
		statusStr := ctx.GetFromArray("utilities.status", anime.GetStatus())

		eb := embed.NewEmbed(ctx.Locale, "utilities.anime.embed").
			WithAuthor("https://cdn.discordapp.com/avatars/743538534589267990/a6c5e905673d041a88b49203d6bc74dd.png?size=2048", "", typeStr, anime.Episodes).
			SetTitle(utils.Fmt("%s | %s", constants.HAPPY, anime.Title.JP)).
			SetDescription(rawSynopsis).
			SetURL(anime.SiteURL).
			SetThumbnail(anime.Cover.ExtraLarge).
			SetImage(anime.Banner).
			SetColor(utils.ToHexNumber(anime.Cover.Color)).
			WithField(strings.Join(anime.GetDirectors(), "\n"), true).
			WithField(strings.Join(anime.GetAnimationStudios(), "\n"), true).
			WithField(anime.GetCreator(), true).
			WithField(sourceStr, true).
			WithField(strings.Join(ctx.GetPrettyGenres(anime.Genres), ", "), true).
			WithField(seasonStr, true).
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

			translatedSynopsis, err := i18n.FromGoogle("en", strings.Split(ctx.ID, "_")[0], rawSynopsis)

			if err == nil {
				eb.SetDescription(translatedSynopsis)
				ctx.EditWithEmbed(msg.ID, eb)
			}

		}

		mal, err := anime.GetBasicFromMAL()

		if err == nil {
			if mal.Score > 0 {
				eb.SetFieldValue(6, utils.Fmt("%.2f", mal.Score))
			}

			if len(mal.Genres) > 0 {
				totalGenres := append(anime.Genres, mal.Genres...)
				eb.SetFieldValue(4, strings.Join(ctx.GetPrettyGenres(totalGenres), ", "))
			}

			ctx.EditWithEmbed(msg.ID, eb)
		}
	},
}
