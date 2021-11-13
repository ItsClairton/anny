package events

import (
	"time"

	"github.com/ItsClairton/Anny/audio"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
)

func OnServerChange(e *gateway.VoiceServerUpdateEvent) {
	time.Sleep(1 * time.Second)

	player := audio.GetPlayer(e.GuildID)
	if player != nil && player.State == audio.PlayingState {
		if err := player.Connection.Speaking(voicegateway.Microphone); err != nil {
			player.Kill(true)
		}
	}
}
