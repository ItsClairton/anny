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

var AnimeCommand = base.Command{
	Name: "anime", Description: "Saber informações básicas sobre um anime",
	Handler: func(ctx *base.CommandContext) {

		if ctx.Args == nil {
			ctx.Reply(Emotes.MIKU_CRY, "VocÊ precisa falar o nome do anime.")
			return
		}

		anime, err := anilist.SearchMediaAsAnime(strings.Join(ctx.Args, " "))

		if err != nil {
			if err.Error() == "Not Found." {
				ctx.Reply(Emotes.MIKU_CRY, "Não encontrei informações sobre esse anime, Desculpa ;(")
			} else {
				ctx.Reply(Emotes.MIKU_CRY, sutils.Fmt("Houve um erro ao obter informações sobre esse anime, desculpe. (%s)", err.Error()))
			}
			return
		}

		launchStr := sutils.Fmt("%s", anime.GetPrettyStartDate())

		if anime.EndDate.Year > 0 && anime.StartDate != anime.EndDate {
			launchStr += sutils.Fmt("\n%s", anime.GetPrettyEndDate())
		}

		if len(anime.GetTrailerURL()) > 0 {
			launchStr = sutils.Fmt("[%s](%s)", launchStr, anime.GetTrailerURL())
		}

		rawSynopsis := sutils.ToMD(anime.Synopsis)

		if err != nil {
			ctx.Reply(Emotes.MIKU_CRY, sutils.Fmt("Um erro ocorreu ao obter a tradução da sinopse. (%s)", err.Error()))
			return
		}

		eb := embed.NewEmbed().
			SetTitle(sutils.Fmt("📺 %s - %s - %d Episódios", anime.Title.JP, anime.GetPrettyFormat(), anime.Episodes)).
			SetDescription(rawSynopsis).
			SetURL(anime.SiteURL).
			SetThumbnail(anime.Cover.ExtraLarge).
			SetImage(anime.Banner).
			AddField("Direção", strings.Join(anime.GetDirectors(), "\n"), true).
			AddField("Estudio", strings.Join(anime.GetAnimationStudios(), "\n"), true).
			AddField("Criador", anime.GetCreator(), true).
			AddField("Status", anime.GetPrettyStatus(), true).
			AddField("Gêneros", strings.Join(anime.Genres, ", "), true).
			AddField("Temporada", anime.GetPrettySeason(), true).
			AddField("Pontuação", "...", true).
			AddField("Data de Lançamento", launchStr, true).
			AddField("Adaptação", anime.GetPrettySource(), true).
			SetFooter("Powered By AniList & MAL", "https://anilist.co/img/icons/favicon-32x32.png")

		hex, err := sutils.ToHexNumber(anime.Cover.Color)

		if err == nil {
			eb.SetColor(hex)
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

		translatedGenres, err := translate.Translate("auto", "pt", strings.Join(anime.Genres, ", "))

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

		score, err := anime.GetScoreFromMAL()

		if err == nil {
			eb.SetField(6, "Pontuação", sutils.Fmt("%.2f", score), true)
			ctx.EditWithEmbed(msg, eb.MessageEmbed)
		}

	},
}
