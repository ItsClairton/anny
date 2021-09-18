package events

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/bwmarrin/discordgo"
)

func InteractionsEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd, exist := discord.GetCommands()[i.ApplicationCommandData().Name]
	if exist {
		go cmd.Handler(&discord.CommandContext{
			InteractionCreate: i,
			Session:           s,
		})
	}

}
