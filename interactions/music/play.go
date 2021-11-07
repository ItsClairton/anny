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
	}, {
		Name:        "embaralhar",
		Description: "Embaralhar as músicas da fila",
		Required:    false,
		Type:        discordgo.ApplicationCommandOptionBoolean,
	}},
	Handler: func(ctx *discord.InteractionContext) error {
		argument := ctx.ApplicationCommandData().Options[0].StringValue()

		shuffle := false
		if len(ctx.ApplicationCommandData().Options) == 2 {
			shuffle = ctx.ApplicationCommandData().Options[1].BoolValue()
		}

		state := ctx.VoiceState()
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
			defer player.AddSong(ctx.Member.User, shuffle, result.Songs...)

			playlist := result.Songs[0].Playlist
			embed.SetColor(0x00D166).SetThumbnail(result.Songs[0].Thumbnail).
				SetDescription("%s Playlist [%s](%s), carregada com sucesso.", emojis.ZeroYeah, playlist.Title, playlist.URL).
				AddField("Criador", playlist.Author, true).
				AddField("Itens", utils.Fmt("%v", len(result.Songs)), true).
				AddField("Duração", utils.FormatTime(playlist.Duration), true)

			return ctx.Edit()
		}

		song := result.Songs[0]
		embed.AddField("Autor", song.Author, true).
			AddField("Duração", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
			AddField("Provedor", song.Provider(), true)

		if !song.IsLoaded() {
			embed.SetDescription("%s Carregando melhores informações de [%s](%s)...", emojis.AnimatedStaff, song.Title, song.URL)
			ctx.Edit()

			song, err = song.Load()
			if err != nil {
				return ctx.SendWithError(err)
			}
		}

		defer player.AddSong(ctx.Member.User, shuffle, song)
		embed.SetColor(0x00D166).
			SetThumbnail(song.Thumbnail).
			SetDescription("%s Música [%s](%s) adicionada com sucesso na fila", emojis.ZeroYeah, song.Title, song.URL)

		return ctx.Edit()
	},
}
