package events

import (
	"context"
	"time"

	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
)

func VoiceServerUpdate(e *gateway.VoiceServerUpdateEvent) {
	defer func() { // Se você usar Connection.Speaking enquanto a conexão está fechada isso causa um Panic
		if r, p := recover(), audio.GetPlayer(e.GuildID); r != nil && p != nil {
			p.Kill(true, emojis.Cry, "Conexão com o servidor de voz muito instável.")
		}
	}()

	time.Sleep(800 * time.Millisecond)
	if p := audio.GetPlayer(e.GuildID); p != nil && p.State == audio.PlayingState {
		logger.DebugF("Mudança de servidor de voz da guilda %d.", e.GuildID)

		if err := p.Connection.Speaking(context.Background(), voicegateway.Microphone); err != nil {
			p.Kill(true, emojis.Cry, "Falha ao enviar dados para o servidor de voz.")
		}
	}
}

func VoiceStateUpdate(e *gateway.VoiceStateUpdateEvent) {
	if e.UserID != base.Me().ID {
		return
	}

	time.Sleep(800 * time.Millisecond)
	if p := audio.GetPlayer(e.GuildID); p != nil && e.ChannelID.IsNull() && p.State != audio.PlayingState {
		p.Kill(true)

		logs, err := base.Session.AuditLog(e.GuildID, api.AuditLogData{ActionType: discord.MemberDisconnect, Limit: 1})
		if err == nil && len(logs.Users) > 0 && time.Since(logs.Entries[0].CreatedAt()) < 5*time.Second {
			base.SendMessage(p.TextID, emojis.Cry, "O vacilão do <@%s> me expulsou do canal de voz, bonk nele ;(", logs.Entries[0].UserID)
		} else {
			base.SendMessage(p.TextID, emojis.Cry, "Quem foi o vacilão que me expulsou do canal de voz? ;(")
		}
	}
}
