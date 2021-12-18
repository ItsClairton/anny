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
	Deffered:    true,
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "nome",
		Description: "nome da música",
	}},
	Handler: func(ctx *core.CommandContext) {
		player, query := music.GetPlayer(ctx.GuildID), ctx.Argument(0).String()

		var replaced bool
		if query == "" {
			if player != nil && player.Current != nil {
				if !strings.Contains(player.Current.Title, "-") {
					query = utils.Fmt("%s - %s", player.Current.Author, player.Current.Title)
					replaced = true
				} else {
					query = player.Current.Title
				}
			} else {
				ctx.Reply(emojis.Cry, "Não há nada tocando no momento.")
				return
			}
		}

		data, err := gonius.SearchSong(query)
		if err != nil {
			if err != gonius.ErrNotFound {
				ctx.Stacktrace(err)
				return
			} else if replaced && player.Current != nil {
				data, err = gonius.SearchSong(player.Current.Title)
			}
		}

		if err != nil {
			if err != gonius.ErrNotFound {
				ctx.Stacktrace(err)
			} else {
				ctx.Reply(emojis.Cry, "Não consegui achar a letra dessa música.")
			}
			return
		}

		lyrics, err := data.Lyrics()
		if err != nil {
			ctx.Stacktrace(err)
			return
		}

		if len(lyrics) > 4096 {
			ctx.Reply(emojis.Cry, "A letra dessa música é muito grande.")
			return
		}

		ctx.Reply(utils.NewEmbed().
			Color(0x9B59B6).
			Thumbnail(data.Thumbnail).
			Title(data.FullTitle).
			Description(lyrics))
	},
}
