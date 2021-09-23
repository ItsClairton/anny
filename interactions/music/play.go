package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/Pauloo27/searchtube"
	"github.com/bwmarrin/discordgo"
)

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
		channel := ctx.GetVoiceChannel()
		if channel == "" {
			ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Você não está conectado em nenhum canal de voz.")
			return
		}
		ctx.SendDeffered()

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			vc, err := ctx.Session.ChannelVoiceJoin(ctx.GuildID, channel, false, true)
			if err != nil {
				ctx.SendError(err)
				return
			}
			player = audio.AddPlayer(audio.NewPlayer(ctx.GuildID, vc))
		}

		result, err := searchtube.Search(ctx.ApplicationCommandData().Options[0].StringValue(), 1)
		if err != nil {
			ctx.SendError(err)
			return
		}

		track, err := audio.GetTrack(result[0].ID, ctx.User)
		if err != nil {
			ctx.SendError(err)
			return
		}
		player.AddQueue(track)
		ctx.EditWithEmote(emojis.PepeArt, "A música **%s** de **%s** foi adicionada com sucesso na fila.", track.Name, track.Author)
	},
}
