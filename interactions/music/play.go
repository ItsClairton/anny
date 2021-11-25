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
			SetColor(0xF0FF00).
			SetDescription("%s Obtendo resultados para sua pesquisa...", emojis.AnimatedStaff)
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
			embed.SetColor(0xF93A2F).SetDescription("%s Não consegui achar essa música, Desculpa ;(", emojis.Cry)
			return ctx.WithEmbed(embed).Edit()
		}

		if result.IsFromPlaylist {
			defer player.AddSong(&ctx.Member.User, shuffle, result.Songs...)

			playlist := result.Songs[0].Playlist
			embed.SetColor(0x00D166).SetThumbnail(result.Songs[0].Thumbnail).
				SetDescription("%s Lista de reprodução [%s](%s) adicionada na fila", emojis.Yeah, playlist.Title, playlist.URL).
				AddField("Criador", playlist.Author, true).
				AddField("Itens", utils.Fmt("%d", len(result.Songs)), true).
				AddField("Duração", utils.FormatTime(playlist.Duration), true)

			return ctx.WithEmbed(embed).Edit()
		}

		song := result.Songs[0]
		embed.AddField("Autor", song.Author, true).
			AddField("Duração", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
			AddField("Provedor", song.Provider(), true)

		if !song.IsLoaded() {
			embed.SetDescription("%s Carregando mais informações de [%s](%s)...", emojis.AnimatedStaff, song.Title, song.URL)
			ctx.WithEmbed(embed).Edit()

			song, err = song.Load()
			if err != nil {
				return ctx.SendError(err)
			}
		}

		defer player.AddSong(&ctx.Member.User, shuffle, song)
		return ctx.WithEmbed(embed.SetColor(0x00D166).
			SetThumbnail(song.Thumbnail).
			SetDescription("%s Música [%s](%s) adicionada na fila", emojis.Yeah, song.Title, song.URL)).Edit()
	},
}
