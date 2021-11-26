package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

var PlayCommand = core.Interaction{
	Name:        "tocar",
	Description: "Tocar uma música, lista de reprodução, ou live",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "argumento",
		Description: "Titulo, ou URL de um vídeo, áudio, ou playlist",
		Required:    true,
	}, &discord.BooleanOption{
		OptionName:  "embaralhar",
		Description: "Embaralhar as músicas na fila caso seja uma playlist",
	}},
	Handler: func(ctx *core.InteractionContext) error {
		query, shuffle := ctx.ArgumentAsString(0), ctx.ArgumentAsBool(1)

		state := ctx.VoiceState()
		if state == nil {
			return ctx.AsEphemeral().Send(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
		}

		embed := core.
			NewEmbed().
			Color(0xF0FF00).
			Description("%s Obtendo resultados para sua pesquisa...", emojis.AnimatedStaff)
		ctx.WithEmbed(embed).Send()

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			player = audio.NewPlayer(ctx.GuildID, ctx.ChannelID, state.ChannelID)
		}

		defer player.Kill(false)
		result, err := audio.FindSong(query)
		if err != nil {
			return ctx.SendError(err)
		}

		if result == nil {
			embed.Color(0xF93A2F).Description("%s Não consegui achar essa música, Desculpa ;(", emojis.Cry)
			return ctx.WithEmbed(embed).Edit()
		}

		if result.IsFromPlaylist {
			defer player.AddSong(&ctx.Member.User, shuffle, result.Songs...)

			playlist := result.Songs[0].Playlist
			embed.Color(0x00D166).Thumbnail(result.Songs[0].Thumbnail).
				Description("%s Lista de reprodução [%s](%s) adicionada na fila", emojis.Yeah, playlist.Title, playlist.URL).
				Field("Criador", playlist.Author, true).
				Field("Itens", utils.Fmt("%d", len(result.Songs)), true).
				Field("Duração", utils.FormatTime(playlist.Duration), true)

			return ctx.WithEmbed(embed).Edit()
		}

		song := result.Songs[0]
		embed.Field("Autor", song.Author, true).
			Field("Duração", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
			Field("Provedor", song.Provider(), true)

		if !song.IsLoaded() {
			embed.Description("%s Carregando mais informações de [%s](%s)...", emojis.AnimatedStaff, song.Title, song.URL)
			ctx.WithEmbed(embed).Edit()

			song, err = song.Load()
			if err != nil {
				return ctx.SendError(err)
			}
		}

		defer player.AddSong(&ctx.Member.User, shuffle, song)
		return ctx.WithEmbed(embed.Color(0x00D166).
			Thumbnail(song.Thumbnail).
			Description("%s Música [%s](%s) adicionada na fila", emojis.Yeah, song.Title, song.URL)).Edit()
	},
}
