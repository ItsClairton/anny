package music

import (
	"time"

	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

var VServerUpdateEvent = core.Event{
	Handler: func(e *gateway.VoiceServerUpdateEvent) {
		time.Sleep(1 * time.Second)

		if player := music.GetPlayer(e.GuildID); player != nil && player.State == music.PlayingState {
			player.Voicy.SendSpeaking()
		}
	},
}

var VStateUpdateEvent = core.Event{
	Handler: func(e *gateway.VoiceStateUpdateEvent) {
		if e.UserID != core.Self.ID {
			return
		}

		if player := music.GetPlayer(e.GuildID); player != nil && e.ChannelID.IsNull() {
			player.Stop(false)

			logs, err := core.State.AuditLog(e.GuildID, api.AuditLogData{ActionType: discord.MemberDisconnect, Limit: 1})
			if err == nil && len(logs.Users) > 0 && time.Since(logs.Entries[0].CreatedAt()) < 5*time.Second {
				player.Send(emojis.Cry, "O vacilão do <@%s> me expulsou do batidão, bonk nele %s", logs.Entries[0].UserID, emojis.AnimatedBonk)
			} else {
				player.Send(emojis.AnimatedBonk, "Quem foi o vacião que me expulsou do batidão? %s", emojis.Cry)
			}
		}
	},
}
