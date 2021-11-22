package events

import (
	"time"

	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func VoiceServerUpdate(e *gateway.VoiceServerUpdateEvent) {
	time.Sleep(1 * time.Second)
	player := audio.GetPlayer(e.GuildID)
	if player != nil && player.State == audio.PlayingState {
		logger.DebugF("Mudança de Região de voz: %d", e.GuildID)

		if err := player.Session.SendSpeaking(); err != nil {
			player.Kill(true, emojis.Cry, "Conexão com o servidor de voz perdida ;(")
			logger.ErrorF("Um erro ocorreu ao enviar pacote de Speaking para o Discord, ID %d: %v", e.GuildID, err)
		}
	}
}

func VoiceStateUpdate(e *gateway.VoiceStateUpdateEvent) {
	if e.UserID != base.Me().ID {
		return
	}

	time.Sleep(500 * time.Millisecond)
	player := audio.GetPlayer(e.GuildID)

	if player != nil && player.State == audio.PlayingState && e.Mute {
		player.Pause()
	}

	if player != nil && player.State == audio.PausedState && !e.Mute {
		player.Resume()
	}

	if player != nil && e.ChannelID.IsNull() {
		player.Kill(true)

		if author := getActionAuthor(e.GuildID, discord.MemberDisconnect); author.IsValid() {
			base.SendMessage(player.TextID, emojis.Cry, "O Vacilão do <@%d> me desconectou do canal de voz, Bonk nele ;(", author)
		} else {
			base.SendMessage(player.TextID, emojis.Cry, "Fui desconectada do canal de voz, Sayonara ;(")
		}
	}
}

func getActionAuthor(guildID discord.GuildID, action discord.AuditLogEvent) discord.UserID {
	logs, err := base.Session.AuditLog(guildID, api.AuditLogData{
		ActionType: action,
		Limit:      1,
	})

	if err == nil && len(logs.Users) > 0 && time.Since(logs.Entries[0].CreatedAt()) < 5*time.Second {
		return logs.Entries[0].UserID
	}

	return 0
}
