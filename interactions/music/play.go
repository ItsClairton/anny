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
	Description: "Tocar uma música, live ou playlist do YouTube",
	Options: []*discordgo.ApplicationCommandOption{{
		Name:        "argumento",
		Description: "Titulo ou link do conteúdo no YouTube",
		Required:    true,
		Type:        discordgo.ApplicationCommandOptionString,
	}},
	Handler: func(ctx *discord.InteractionContext) {
		argument := ctx.ApplicationCommandData().Options[0].StringValue()
		voiceID := ctx.GetVoiceChannel()
		if voiceID == "" {
			ctx.SendEphemeral(emojis.MikuCry, "Você precisa estar conectado a um canal de voz para utilizar esse comando.")
			return
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			player = audio.NewPlayer(ctx.GuildID, ctx.ChannelID, voiceID)
		}
		player.Lock()
		defer player.Unlock()

		embed := discord.NewEmbed().
			SetColor(0xF8C300).
			SetDescription("%s Obtendo melhores resultados para sua pesquisa...", emojis.AnimatedStaff)
		go ctx.SendEmbed(embed.Build())

		result, err := audio.FindSong(argument)
		if err != nil {
			ctx.SendError(err)
			go player.Kill(false)
			return
		}

		if result == nil {
			ctx.SendEmbed(embed.
				SetColor(0xF93A2F).
				SetDescription("%s Não consegui achar essa música, Desculpa ;(", emojis.MikuCry).
				Build())

			go player.Kill(false)
			return
		}

		if result.IsFromSearch {
			if len(result.Songs) == 1 {
				go loadAndAdd(ctx, player, result.Songs[0])
				return
			}

			description := ""
			response := discord.NewResponse().WithEmbed(embed)

			for i := range result.Songs {
				song := result.Songs[i]
				description += utils.Fmt("\n%s [%s](%s) de %s", emojis.GetNumberAsEmoji(i+1), song.Title, song.URL, song.Author)

				response.WithButton(discord.Button{
					UserID:  ctx.Member.User.ID,
					Style:   discordgo.SecondaryButton,
					Label:   utils.Is(song.IsLive, "LIVE", utils.FormatTime(song.Duration)),
					Emoji:   emojis.GetNumberAsEmoji(i + 1),
					Once:    true,
					Delayed: true,
					OnClick: func(_ *discord.InteractionContext) {
						loadAndAdd(ctx, player, song)
					},
				})
			}

			embed.SetColor(0x0099e1).SetDescription(description)
			ctx.SendResponse(response)
			return
		}

		if result.IsFromPlaylist {
			go func() {
				player.AddSong(ctx.Member.User, result.Songs...)

				playlist := result.Songs[0].Playlist
				ctx.SendEmbed(embed.
					SetColor(0x00D166).
					SetThumbnail(result.Songs[0].Thumbnail).
					SetDescription(utils.Fmt("%s A lista de reprodução [%s](%s) foi carregada com sucesso.", emojis.ZeroYeah, playlist.Title, playlist.URL)).
					AddField("Autor", playlist.Author, true).
					AddField("Itens", utils.Fmt("%v", len(result.Songs)), true).
					AddField("Duração", utils.FormatTime(playlist.Duration), true).Build())

			}()
			return
		}

		go loadAndAdd(ctx, player, result.Songs[0])
	},
}

func loadAndAdd(ctx *discord.InteractionContext, player *audio.Player, song *audio.Song) {
	player.Lock()
	embed := discord.NewEmbed().SetColor(0xF8C300).
		AddField("Autor", song.Author, true).
		AddField("Duração", utils.FormatTime(song.Duration), true).
		AddField("Provedor", song.Provider.Name(), true)

	if song.StreamingURL == "" {
		go ctx.SendEmbed(embed.SetDescription(utils.Fmt("%s Tentando descriptografar [%s](%s)...", emojis.AnimatedStaff, song.Title, song.URL)).Build())

		var err error
		song, err = song.Provider.GetInfo(song)
		if err != nil {
			ctx.SendError(err)
			player.Unlock()
			player.Kill(false)
			return
		}
	}
	player.Unlock()
	player.AddSong(ctx.Member.User, song)

	ctx.SendEmbed(embed.SetColor(0x00D166).
		SetThumbnail(song.Thumbnail).
		SetDescription(utils.Fmt("%s [%s](%s) foi adicionado com sucesso na fila", emojis.ZeroYeah, song.Title, song.URL)).
		Build())
}
