package events

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/bwmarrin/discordgo"
)

func InteractionsEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		cmd, exist := discord.GetInteractions()[i.ApplicationCommandData().Name]
		if exist {
			go cmd.Handler(&discord.InteractionContext{
				InteractionCreate: i,
				Session:           s,
			})
		}
		return
	}

	if i.Type == discordgo.InteractionMessageComponent {
		button := discord.GetButton(i.MessageComponentData().CustomID)
		if button != nil {
			go button.OnClick(&discord.InteractionContext{
				InteractionCreate: i,
				Session:           s,
			})

			if button.Once {
				discord.UnregisterButton(button.ID)
			}
		} else {
			discord.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Desculpe, essa interação já expirou.",
					Flags:   1 << 6,
				},
			})
		}
	}
}
