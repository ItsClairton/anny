package music

import (
	"net/url"
	"strings"

	"github.com/ItsClairton/Anny/core"
	"github.com/tidwall/gjson"

	music "github.com/ItsClairton/Anny/music/audio"
	providers "github.com/ItsClairton/Anny/music/providers"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

var PlayCommand = core.Command{
	Name:        "tocar",
	Description: "Sistema de músicas",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:   "musica",
		Description:  "Nome, ou URL de uma música ou playlist",
		Required:     true,
		Autocomplete: true,
	}, &discord.BooleanOption{
		OptionName:  "embaralhar",
		Description: "Embaralhar as músicas da fila",
	}},
	Handler: func(ctx *core.CommandContext) {
		query, shuffle := ctx.Argument(0).String(), ctx.Argument(1).Bool()

		state := ctx.VoiceState()
		if state == nil {
			ctx.Ephemeral().Reply(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
			return
		}

		embed := utils.NewEmbed().Color(0xF0FF00).Description("%s Obtendo resultados para sua pesquisa...", emojis.AnimatedStaff)
		ctx.Reply(embed)

		player := music.GetOrCreatePlayer(ctx.GuildID, ctx.ChannelID, state.ChannelID)
		defer checkIdle(player)

		result, err := providers.FindSong(query, true)
		if err != nil {
			ctx.Stacktrace(err)
			return
		}

		if result == nil {
			ctx.Reply(embed.Color(0xF93A2F).Description("%s Não consegui encontrar essa música.", emojis.Cry))
			return
		}

		if result.Playlist != nil {
			player.AddSong(ctx.Sender(), shuffle, result.Songs...)

			ctx.Reply(embed.Color(0x00D166).
				Description("%s Lista de reprodução [%s](%s) adicionada na fila", emojis.Yeah, result.Playlist.Title, result.Playlist.URL).
				Field("Criador", result.Playlist.Author, true).
				Field("Músicas", len(result.Songs), true).
				Field("Duração", utils.FormatTime(result.Playlist.Duration), true))
			return
		}

		song := result.Songs[0]
		embed.Thumbnail(song.Thumbnail).
			Field("Autor", song.Author, true).
			Field("Duração", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
			Field("Provedor", song.Provider(), true)

		if !song.IsLoaded() {
			go ctx.Reply(embed.Description("%s Carregando [%s](%s)", emojis.AnimatedStaff, song.Title, song.URL))

			if err := song.Load(); err != nil {
				ctx.Stacktrace(err)
				return
			}
		}

		player.AddSong(ctx.Sender(), shuffle, song)
		ctx.Reply(embed.
			Color(0x00D166).
			Thumbnail(song.Thumbnail).
			Description("%s Música [%s](%s) adicionada na fila", emojis.Yeah, song.Title, song.URL))
	},
	AutoCompleteHandler: func(ctx *core.AutoCompleteContext) api.AutocompleteStringChoices {
		query := strings.ReplaceAll(ctx.Data.Options[0].Value.String(), "\"", "")
		if strings.TrimSpace(query) == "" {
			return api.AutocompleteStringChoices{}
		}

		if providers.FindByInput(query, false) != nil {
			return api.AutocompleteStringChoices{{
				Name:  query,
				Value: query,
			}}
		}

		data, err := utils.FromWebString("http://suggestqueries.google.com/complete/search?client=youtube&ds=yt&client=chrome&q=" + url.QueryEscape(query))
		if err != nil {
			panic(err)
		}

		result := api.AutocompleteStringChoices{}
		for _, entry := range gjson.Get(data, "1").Array() {
			choice := entry.String()

			result = append(result, discord.StringChoice{
				Name:  choice,
				Value: choice,
			})
		}

		return result
	},
}

func checkIdle(player *music.Player) {
	if player.State != music.StoppedState || len(player.Queue) != 0 {
		return
	}

	player.Stop(true)
}
