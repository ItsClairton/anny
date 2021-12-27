package music

import (
	"strings"

	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/gonius"
	"github.com/diamondburned/arikawa/v3/discord"
)

var LyricsCommand = core.Command{
	Name:        "letra",
	Description: "Mostra a letra da música",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "nome",
		Description: "nome da música",
	}},
	Handler: func(ctx *core.CommandContext) {
		player, query, replaced := music.GetPlayer(ctx.GuildID), ctx.Argument(0).String(), false
		if query == "" {
			if player != nil && player.Current != nil {
				if strings.Contains(player.Current.Title, "-") {
					query = player.Current.Title
				} else {
					query, replaced = utils.Fmt("%s - %s", player.Current.Author, player.Current.Title), true
				}
			} else {
				ctx.Reply(emojis.Cry, "Não há nada tocando no momento, e você não passou nenhuam música para obter a letra.")
				return
			}
		}

		embed := utils.NewEmbed().Color(0xF0FF00).Description("%s Obtendo resultados...", emojis.AnimatedStaff)
		ctx.Reply(embed)

		result, err := gonius.SearchSong(query)
		if err != nil {
			if err != gonius.ErrNotFound {
				ctx.Stacktrace(err)
				return
			}

			if replaced && player.Current != nil {
				result, err = gonius.SearchSong(player.Current.Title)
			}
		}

		if err != nil {
			if err == gonius.ErrNotFound {
				ctx.Reply(embed.Color(0xF93A2F).Description("%s Não consegui achar a letra dessa música.", emojis.Cry))
			} else {
				ctx.Stacktrace(err)
			}
			return
		}

		entry := result[0]
		ctx.Reply(embed.Description("%s Carregando letra...", emojis.AnimatedStaff))

		lyrics, err := entry.Lyrics()
		if err != nil {
			ctx.Stacktrace(err)
			return
		}

		if len(lyrics) > 4096 {
			ctx.Reply(embed.Color(0xF93A2F).Description("%s A letra dessa música é muito grande.", emojis.Cry))
			return
		}

		ctx.Reply(embed.
			Color(0x0099E1).
			Author(entry.PrimaryArtist.Name, entry.PrimaryArtist.Image).
			URL(entry.URL).
			Title(entry.Title).
			Description(lyrics).
			Thumbnail(entry.Image))
	},
}
