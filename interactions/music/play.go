package music

import (
	"regexp"

	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/Pauloo27/searchtube"
	"github.com/bwmarrin/discordgo"
)

var videoRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtu\.?be(?:\.com)?\/?.*(?:watch|embed)?(?:.*v=|v\/|\/)([\w\-_]+)\&?`)
var playlistRegex = regexp.MustCompile(`^https?:\/\/(www.youtube.com|youtube.com)\/playlist(.*)$`)

var PlayCommand = discord.Interaction{
	Name:        "tocar",
	Description: "Toca algum vídeo do YouTube em um canal de voz",
	Options: []*discordgo.ApplicationCommandOption{{
		Name:        "vídeo",
		Description: "Titulo ou link de um vídeo do YouTube",
		Type:        discordgo.ApplicationCommandOptionString,
		Required:    true,
	}},
	Handler: func(ctx *discord.InteractionContext) {
		voiceId := ctx.GetVoiceChannel()
		if voiceId == "" {
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
			return
		}
		ctx.SendDeffered(false)

		embed := discord.NewEmbed().
			SetDescription(utils.Fmt("%s Tentando se conectar ao canal...", emojis.AnimatedStaff)).
			SetColor(0xF8C300)
		response := discord.NewResponse().WithEmbed(embed.Build())
		content := ctx.ApplicationCommandData().Options[0].StringValue()

		if playlistRegex.MatchString(content) {
			ctx.EditResponse(response)

			player, err := audio.GetOrCreatePlayer(ctx.Session, ctx.GuildID, ctx.ChannelID, voiceId)
			if err != nil {
				embed.SetColor(0xF93A2F).
					SetDescription(utils.Fmt("%s Um erro ocorreu ao tentar se conectar ao canal: `%s`", emojis.MikuCry, err.Error()))
				ctx.EditResponse(response)
				audio.RemovePlayer(player, false)
				return
			}

			embed.SetDescription(utils.Fmt("%s Obtendo informações detalhadas sobre a Playlist...", emojis.AnimatedStaff))
			ctx.EditResponse(response)

			tracks, time, err := audio.GetPlaylist(content, ctx.Member.User)
			if err != nil {
				embed.SetColor(0xF93A2F).
					SetDescription(utils.Fmt("%s Um erro ocorreu ao carregar essa Playlist: `%s`", emojis.MikuCry, err.Error()))
				ctx.EditResponse(response)
				audio.RemovePlayer(player, false)
				return
			}
			player.AddPlaylist(tracks)
			embed.
				SetDescription(utils.Fmt("A Playlist [%s](%s) foi adicionada na fila", tracks[0].Playlist.Title, utils.Fmt("https://www.youtube.com/playlist?list=%s", tracks[0].Playlist.ID))).
				SetColor(0x00D166).
				AddField("Autor", tracks[0].Playlist.Author, true).
				AddField("Itens", utils.Fmt("%d", len(tracks)), true).
				AddField("Duração", time.String(), true).
				SetImage(tracks[0].ThumbnailUrl)

			ctx.EditResponse(response)
			return
		}

		if videoRegex.MatchString(content) {
			ctx.EditResponse(response)

			player, err := audio.GetOrCreatePlayer(ctx.Session, ctx.GuildID, ctx.ChannelID, voiceId)
			if err != nil {
				embed.SetColor(0xF93A2F).SetDescription(utils.Fmt("%s Um erro ocorreu ao tentar se conectar ao canal: `%s`", emojis.MikuCry, err.Error()))
				ctx.EditResponse(response)
				return
			}

			embed.SetDescription(utils.Fmt("%s Tentando decodificar á música...", emojis.AnimatedStaff))
			ctx.EditResponse(response)
			player.Lock()

			track, err := audio.GetTrack(content, ctx.Member.User)
			if err != nil {
				embed.SetColor(0xF93A2F).
					SetDescription(utils.Fmt("%s Um erro ocorreu ao decodificar essa música: `%s`", emojis.MikuCry, err.Error()))
				ctx.EditResponse(response)

				player.Unlock()
				audio.RemovePlayer(player, false)
				return
			}

			player.Unlock()
			player.AddTrack(track)

			embed.SetColor(0x00D166).
				SetDescription(utils.Fmt("A música [%s](%s) foi adicionada com sucesso na fila", track.Title, track.URL)).
				SetImage(track.ThumbnailUrl).
				AddField("Autor", track.Author, true).
				AddField("Duração", utils.ToDisplayTime(track.Duration.Seconds()), true)

			ctx.EditResponse(response)
			return
		}

		result, err := searchtube.Search(content, 5)
		if err != nil {
			ctx.SendError(err)
			return
		}

		if len(result) <= 0 {
			ctx.EditWithEmote(emojis.MikuCry, "Não foi possível encontrar essa música.")
			return
		}
		if len(result) == 1 {
			handleResult(ctx, voiceId, result[0], ctx.Member.User)
			return
		}

		var resultText string

		for i := range result {
			entry := result[i]
			resultText += utils.Fmt("%s %s de **%s**\n", emojis.GetNumberAsEmoji(i+1), entry.Title, entry.Uploader)
			response.WithButton(discord.Button{
				Label: utils.Is(entry.Live, "LIVE", entry.RawDuration),
				Once:  true,
				Emoji: emojis.GetNumberAsEmoji(i + 1),
				Style: discordgo.SecondaryButton,
				OnClick: func(btx *discord.InteractionContext) {
					handleResult(ctx, voiceId, entry, btx.Member.User)
				},
			})
		}

		ctx.EditResponse(response.ClearEmbeds().WithEmbed(discord.NewEmbed().
			SetColor(0x00D166).
			SetDescription(resultText).Build()))

	},
}

func handleResult(ctx *discord.InteractionContext, voiceId string, entry *searchtube.SearchResult, user *discordgo.User) {
	embed := discord.NewEmbed().
		SetDescription(utils.Fmt("%s Tentando se conectar ao canal...", emojis.AnimatedStaff)).
		SetColor(0xF8C300)

	response := discord.NewResponse().ClearComponents().WithEmbed(embed.Build())
	ctx.EditResponse(response)

	player, err := audio.GetOrCreatePlayer(ctx.Session, ctx.GuildID, ctx.ChannelID, voiceId)
	if err != nil {
		embed.SetColor(0xF93A2F).SetDescription(utils.Fmt("%s Um erro ocorreu ao tentar se conectar ao canal: `%s`", emojis.MikuCry, err.Error()))
		ctx.EditResponse(response)
		return
	}

	thumbnailUrl := utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", entry.ID)
	embed.SetDescription(utils.Fmt("%s Decodificando: [%s](%s)", emojis.AnimatedStaff, entry.Title, entry.URL)).
		SetImage(thumbnailUrl).
		AddField("Autor", entry.Uploader, true).
		AddField("Duração", entry.RawDuration, true)
	ctx.EditResponse(response)

	player.Lock()
	track, err := audio.GetTrack(entry.ID, user)
	if err != nil {
		embed.SetColor(0xF93A2F).
			SetDescription(utils.Fmt("%s Um erro ocorreu ao decodificar essa música: `%s`", emojis.MikuCry, err.Error()))
		ctx.EditResponse(response)

		player.Unlock()
		audio.RemovePlayer(player, false)
		return
	}

	player.Unlock()
	player.AddTrack(track)
	embed.SetColor(0x00D166).
		SetImage(track.ThumbnailUrl).
		SetDescription(utils.Fmt("A música [%s](%s) foi adicionada com sucesso na fila", entry.Title, entry.URL))
	ctx.EditResponse(response)
}
