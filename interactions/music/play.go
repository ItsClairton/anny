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

var regex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtu\.?be(?:\.com)?\/?.*(?:watch|embed)?(?:.*v=|v\/|\/)([\w\-_]+)\&?`)

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
		ctx.SendDeffered(true)

		content := ctx.ApplicationCommandData().Options[0].StringValue()
		if regex.MatchString(content) {
			player, err := audio.GetOrCreatePlayer(ctx.Session, ctx.GuildID, ctx.ChannelID, voiceId)
			if err != nil {
				ctx.SendError(err)
				return
			}
			player.Lock()
			track, err := audio.GetTrack(content, ctx.User)
			if err != nil {
				ctx.SendError(err)
				player.Unlock()
				audio.RemovePlayer(player, false)
				return
			}

			player.Unlock()
			player.AddQueue(track)
			if player.GetState() != audio.StoppedState {
				ctx.EditWithEmote(emojis.PepeArt, "A música **%s** de **%s** foi adicionada com sucesso na posição **%d** da fila.", track.Title, track.Author, len(player.GetQueue()))
			} else {
				ctx.EditWithEmote(emojis.PepeArt, "Tocando agora a música **%s** de **%s**.", track.Title, track.Author)
			}
		} else {
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
				handleResult(ctx, voiceId, result[0])
				return
			}

			response := discord.NewResponse()
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
						ctx.EditResponse(response.ClearComponents())
						handleResult(btx, voiceId, entry)
					},
				})
			}

			embed := discord.NewEmbed().
				SetColor(0x00D166).
				SetDescription(resultText)

			ctx.EditResponse(response.WithEmbed(embed.Build()))
		}
	},
}

func handleResult(ctx *discord.InteractionContext, voiceId string, entry *searchtube.SearchResult) {
	player, err := audio.GetOrCreatePlayer(ctx.Session, ctx.GuildID, ctx.ChannelID, voiceId)
	if err != nil {
		ctx.SendError(err)
		return
	}

	ctx.ReplyWithEmote(emojis.AnimatedStaff, "Tentando decodificar do YouTube **%s** de **%s**...", entry.Title, entry.Uploader)
	player.Lock()
	track, err := audio.GetTrack(entry.ID, ctx.User)
	if err != nil {
		ctx.EditWithEmote(emojis.MikuCry, "Um erro ocorreu ao decodificar essa música. (`%s`)", err.Error())
		player.Unlock()
		audio.RemovePlayer(player, false)
		return
	}
	player.Unlock()
	player.AddQueue(track)
	if player.GetState() != audio.StoppedState {
		ctx.EditWithEmote(emojis.PepeArt, "A música **%s** de **%s** foi adicionada com sucesso na posição **%d** da fila.", track.Title, track.Author, len(player.GetQueue()))
	} else {
		ctx.EditWithEmote(emojis.PepeArt, "Tocando agora a música **%s** de **%s**.", track.Title, track.Author)
	}
}
