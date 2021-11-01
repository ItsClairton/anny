package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

var PlayCommand = discord.Interaction{
	Name:        "tocar",
	Description: "Tocar uma música, lista de reprodução, ou live",
	Options: []*discordgo.ApplicationCommandOption{{
		Name:        "argumento",
		Description: "Titulo ou link do conteúdo no YouTube",
		Required:    true,
		Type:        discordgo.ApplicationCommandOptionString,
	}},
	Handler: func(ctx *discord.InteractionContext) error {
		argument, state := ctx.ApplicationCommandData().Options[0].StringValue(), ctx.VoiceState()
		if state == nil {
			return ctx.AsEphemeral().Send(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			player = audio.NewPlayer(ctx.GuildID, ctx.ChannelID, state.ChannelID)
		}

		embed := discord.
			NewEmbed().
			SetColor(0xF8C300).
			SetDescription("%s Obtendo melhores resultados para sua pesquisa...", emojis.AnimatedStaff)
		ctx.WithEmbed(embed).Send()
		defer player.Kill(false)

		result, err := audio.FindSong(argument)
		if err != nil {
			return ctx.SendWithError(err)
		}

		if result == nil {
			embed.SetColor(0xF93A2F).SetDescription("%s Não consegui achar essa música, Desculpa ;(", emojis.MikuCry)
			return ctx.Edit()
		}

		if result.IsFromPlaylist {
			playlist := result.Songs[0].Playlist
			embed.SetColor(0x00D166).SetThumbnail(result.Songs[0].Thumbnail).
				SetDescription("%s Playlist [%s](%s), carregada com sucesso.", emojis.ZeroYeah, playlist.Title, playlist.URL).
				AddField("Criador", playlist.Author, true).
				AddField("Itens", utils.Fmt("%v", len(result.Songs)), true).
				AddField("Duração", utils.FormatTime(playlist.Duration), true)
			ctx.Edit()

			player.AddSong(ctx.Member.User, result.Songs...)
			return nil
		}

		song := result.Songs[0]
		embed.AddField("Autor", song.Author, true).
			AddField("Duração", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
			AddField("Provedor", song.Provider.Name(), true)

		if song.StreamingURL == "" {
			embed.SetDescription("%s Carregando melhores informações de [%s](%s)...", emojis.AnimatedStaff, song.Title, song.URL)
			ctx.Edit()

			song, err = song.GetMoreInfo()
			if err != nil {
				return ctx.SendWithError(err)
			}
		}
		embed.SetColor(0x00D166).
			SetThumbnail(song.Thumbnail).
			SetDescription("%s Música [%s](%s) adicionada com sucesso na fila", emojis.ZeroYeah, song.Title, song.URL)
		ctx.Edit()

		player.AddSong(ctx.Member.User, song)
		return nil
	},
}
