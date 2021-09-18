package events

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func ReadyEvent(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateListeningStatus(os.Getenv("DISCORD_STATUS"))
}
