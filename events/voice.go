package events

import (
	"time"

	"github.com/ItsClairton/Anny/audio"
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
)

func VoiceServerUpdate(e *gateway.VoiceServerUpdateEvent) {
	time.Sleep(1 * time.Second)
	player := audio.GetPlayer(e.GuildID)
	if player != nil && player.State == audio.PlayingState {
		logger.DebugF("Mudança de Região de voz: %d", e.GuildID)

		if err := player.Connection.Speaking(voicegateway.Microphone); err != nil {
			logger.ErrorF("Um erro ocorreu ao enviar pacote de Speaking para o Discord, ID %d: %v", e.GuildID, err)
		}
	}
}

func VoiceStateUpdate(e *gateway.VoiceStateUpdateEvent) {
	if e.UserID != base.Me().ID {
		return
	}

	player := audio.GetPlayer(e.GuildID)
	if player != nil && e.ChannelID.IsNull() {
		player.Kill(true, emojis.Cry, "Fui desconectada do canal de voz, Sayonara ;(")
	}
}
