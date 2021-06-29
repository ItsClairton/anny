package listeners

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func ReadyListener(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateGameStatus(0, os.Getenv("DISCORD_STATUS"))
}
