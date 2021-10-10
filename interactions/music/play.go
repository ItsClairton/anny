package music

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

var PlayCommand = discord.Interaction{
	Name:        "tocar",
	Description: "Tocar uma música, live ou playlist do YouTube",
	Delayed:     true,
	Options: []*discordgo.ApplicationCommandOption{{
		Name:        "argumento",
		Description: "Titulo ou link do conteúdo no YouTube",
		Required:    true,
		Type:        discordgo.ApplicationCommandOptionString,
	}},
	Handler: func(ctx *discord.InteractionContext) {
		voiceID := ctx.GetVoiceChannel()
		if voiceID == "" {
			ctx.Send(emojis.MikuCry, "Você precisa estar conectado a um canal de voz para utilizar esse comando.")
			return
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			connection := ctx.Session.VoiceConnections[ctx.GuildID]

			if connection == nil {
				connection, err := ctx.Session.ChannelVoiceJoin(ctx.GuildID, voiceID, false, true)
				if err != nil {
					ctx.SendError(err)
					return
				}
				player = audio.NewPlayer(ctx.GuildID, ctx.ChannelID, connection.ChannelID, connection)
			} else {
				player = audio.NewPlayer(ctx.GuildID, ctx.ChannelID, connection.ChannelID, connection)
			}
		}

	},
}
