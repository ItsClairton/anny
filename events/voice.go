package events

import (
	"github.com/ItsClairton/Anny/audio"
	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, state *discordgo.VoiceStateUpdate) {
	if state.BeforeUpdate != nil && state.ChannelID == "" && state.UserID == s.State.User.ID {
		player, connection := audio.GetPlayer(state.GuildID), s.VoiceConnections[state.GuildID]

		if connection != nil && player != nil {
			player.Kill(true)
		}
	}
}
