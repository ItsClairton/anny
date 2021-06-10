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

var AnimeCommand = base.Command{
	Name: "anime",
	Handler: func(ctx *base.CommandContext) {

		if ctx.Args == nil {
			ctx.ReplyWithUsage("<nome de um anime>")
			return
		}

		anime, err := anilist.SearchMediaAsAnime(strings.Join(ctx.Args, " "))

		if err != nil {
			if err.Error() == "Not Found." {
				ctx.Reply(Emotes.MIKU_CRY, "anime.anime.not-found")
			} else {
				ctx.ReplyWithError(err)
			}
			return
		}

		launchStr := sutils.Fmt("%s", date.ToPrettyDate(ctx.Locale, &anime.StartDate))

		if anime.Status == "NOT_YET_RELEASED" {
			launchStr = ctx.Locale.GetString("prevDate", launchStr)
		}

		if anime.EndDate.Year > 0 && anime.StartDate != anime.EndDate {
			launchStr = ctx.Locale.GetString("untilDate", launchStr, date.ToPrettyDate(ctx.Locale, &anime.EndDate))
		}

		hasTrailer := len(anime.GetTrailerURL()) > 0

		if hasTrailer {
			launchStr = sutils.Fmt("[%s](%s)", launchStr, anime.GetTrailerURL())
		}

		rawSynopsis := sutils.ToMD(anime.Synopsis)

		if err != nil {
			ctx.ReplyWithError(err)
			return
		}

		typeStr := ctx.Locale.GetFromArray("anime.type", anime.GetType())
		sourceStr := ctx.Locale.GetFromArray("anime.source", anime.GetSource())
		seasonStr := ctx.Locale.GetFromArray("anime.season", anime.GetSeason())
		statusStr := ctx.Locale.GetFromArray("anime.status", anime.GetStatus())

		eb := embed.NewEmbed(ctx.Locale, "anime.anime.embed").
			WithAuthor("https://cdn.discordapp.com/avatars/743538534589267990/a6c5e905673d041a88b49203d6bc74dd.png?size=2048", "", typeStr, anime.Episodes).
			SetTitle(sutils.Fmt("%s | %s", Emotes.HAPPY, anime.Title.JP)).
			SetDescription(rawSynopsis).
			SetURL(anime.SiteURL).
			SetThumbnail(anime.Cover.ExtraLarge).
			SetImage(anime.Banner).
			SetColor(sutils.ToHexNumber(anime.Cover.Color)).
			WithField(strings.Join(anime.GetDirectors(), "\n"), true).
			WithField(strings.Join(anime.GetAnimationStudios(), "\n"), true).
			WithField(anime.GetCreator(), true).
			WithField(sourceStr, true).
			WithField(strings.Join(ctx.Locale.GetPrettyGenres(anime.Genres), ", "), true).
			WithField(seasonStr, true).
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

			translatedSynopsis, err := i18n.FromGoogle("en", strings.Split(ctx.Locale.ID, "_")[0], rawSynopsis)

			if err == nil {
				eb.SetDescription(translatedSynopsis)
				ctx.EditWithResponse(msg.ID, response)
			}

		}

		mal, err := anime.GetBasicFromMAL()

		if err == nil {
			if mal.Score > 0 {
				eb.SetFieldValue(6, sutils.Fmt("%.2f", mal.Score))
			}

			if len(mal.Genres) > 0 {
				totalGenres := append(anime.Genres, mal.Genres...)
				eb.SetFieldValue(4, strings.Join(ctx.Locale.GetPrettyGenres(totalGenres), ", "))
			}

			ctx.EditWithResponse(msg.ID, response)
		}
	},
}
