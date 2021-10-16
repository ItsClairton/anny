package music

import (
	"regexp"
	"time"

	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/providers"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

var regex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

var PlayCommand = discord.Interaction{
	Name:        "tocar",
	Description: "Tocar uma música, live ou playlist do YouTube",
	Deffered:    true,
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
			ctx.Send(emojis.MikuCry, "Você precisa estar conectado a um canal de voz para utilizar esse comando.")
			return
		}

		player := audio.GetPlayer(ctx.GuildID)
		if player == nil {
			player = audio.NewPlayer(ctx.GuildID, ctx.ChannelID, "", nil)
			audio.AddPlayer(player)
			go func() {
				connection, err := ctx.Session.ChannelVoiceJoin(ctx.GuildID, voiceID, false, true)
				if err != nil {
					if _, ok := ctx.Session.VoiceConnections[ctx.GuildID]; ok {
						connection = ctx.Session.VoiceConnections[ctx.GuildID]
					} else {
						audio.RemovePlayer(player, true)
						ctx.SendError(err)
						return
					}
				}
				player.UpdateVoice(connection.ChannelID, connection)
			}()
		}
		player.Lock()

		if !regex.MatchString(argument) {
			argument = "ytsearch:" + argument
		}

		info, err := providers.FindSong(argument)
		if err != nil {
			ctx.SendError(err)
			player.Unlock()
			audio.RemovePlayer(player, false)
			return
		}

		player.Unlock()
		player.AddTrack(&audio.Track{
			Song:      info,
			Requester: ctx.Member.User,
			Time:      time.Now(),
		})

		embed := discord.NewEmbed().SetColor(0x00D166).
			SetDescription(utils.Fmt("%s O conteúdo [%s](%s) foi adicionado com sucesso na fila", emojis.ZeroYeah, info.Title, info.PageURL)).
			SetImage(info.ThumbnailURL).
			AddField("Autor", info.Uploader, true).
			AddField("Duração", info.Duration(), true).
			AddField("Provedor", info.DisplayProvider(), true).
			Build()

		ctx.SendEmbed(embed)
	},
}
