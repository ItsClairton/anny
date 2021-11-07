package events

import (
	"time"

	"github.com/ItsClairton/Anny/audio"
	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, state *discordgo.VoiceStateUpdate) {
	if state.BeforeUpdate != nil && state.ChannelID == "" && state.UserID == s.State.User.ID {
		time.Sleep(500 * time.Millisecond) // Dar tempo do DiscordGo deletar a conexão do Map

		player, connection := audio.GetPlayer(state.GuildID), s.VoiceConnections[state.GuildID]
		if connection == nil && player != nil {
			player.Kill(true)
		}
	}
}
