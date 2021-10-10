package events

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/bwmarrin/discordgo"
)

func InteractionsEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand: // Slash Commands, Context Menu's
		cmd, exist := discord.GetInteractions()[i.ApplicationCommandData().Name]
		if exist {
			if cmd.Delayed {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: 5})
				go cmd.Handler(&discord.InteractionContext{Session: s, AlreadySended: true, InteractionCreate: i})
			} else {
				go cmd.Handler(&discord.InteractionContext{Session: s, ResponseType: 4, InteractionCreate: i})
			}
		}

	case discordgo.InteractionMessageComponent: // Components
		button := discord.GetButton(i.MessageComponentData().CustomID)
		if button != nil {
			if button.Delayed {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: 6})
				go button.OnClick(&discord.InteractionContext{Session: s, AlreadySended: true, InteractionCreate: i})
			} else {
				go button.OnClick(&discord.InteractionContext{Session: s, ResponseType: 4, InteractionCreate: i})
			}
			if button.Once {
				discord.UnregisterButton(button.ID)
			}
		}
	}
}
