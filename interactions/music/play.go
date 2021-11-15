package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

var PlayCommand = base.Interaction{
	Name:        "tocar",
	Description: "Tocar uma música, lista de reprodução, ou live",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "argumento",
		Description: "Titulo, ou URL de um vídeo, playlist ou áudio",
		Required:    true,
	}, &discord.BooleanOption{
		OptionName:  "embaralhar",
		Description: "Embaralhar as músicas na fila caso seja uma playlist",
	}},
	Handler: func(ctx *base.InteractionContext) error {
		query, shuffle := ctx.ArgumentAsString(0), ctx.ArgumentAsBool(1)

		state := ctx.VoiceState()
		if state == nil {
			return ctx.AsEphemeral().Send(emojis.Cry, "Você não está conectado em nenhum canal de voz.")
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			player = audio.NewPlayer(ctx.GuildID, ctx.ChannelID, state.ChannelID)
		}

		embed := base.
			NewEmbed().
			SetColor(0xF8C300).
			SetDescription("%s Obtendo melhores resultados para sua pesquisa...", emojis.AnimatedStaff)
		ctx.WithEmbed(embed).Send()
		defer player.Kill(false, "", "")

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
				SetDescription("%s Playlist [%s](%s), carregada com sucesso.", emojis.Yeah, playlist.Title, playlist.URL).
				AddField("Criador", playlist.Author, true).
				AddField("Itens", utils.Fmt("%v", len(result.Songs)), true).
				AddField("Duração", utils.FormatTime(playlist.Duration), true)

			return ctx.WithEmbed(embed).Edit()
		}

		song := result.Songs[0]
		embed.AddField("Autor", song.Author, true).
			AddField("Duração", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
			AddField("Provedor", song.Provider(), true)

		if !song.IsLoaded() {
			embed.SetDescription("%s Carregando melhores informações de [%s](%s)...", emojis.AnimatedStaff, song.Title, song.URL)
			ctx.WithEmbed(embed).Edit()

			song, err = song.Load()
			if err != nil {
				return ctx.SendError(err)
			}
		}

		defer player.AddSong(&ctx.Member.User, shuffle, song)
		embed.SetColor(0x00D166).
			SetThumbnail(song.Thumbnail).
			SetDescription("%s Música [%s](%s) adicionada com sucesso na fila", emojis.Yeah, song.Title, song.URL)

		return ctx.WithEmbed(embed).Edit()
	},
}
