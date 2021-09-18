package discord

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/bwmarrin/discordgo"
)

type Category struct {
	Name, Emote string
	Commands    []*Command
}

type Command struct {
	Name, Description string
	Type              discordgo.ApplicationCommandType
	Options           []*discordgo.ApplicationCommandOption
	Category          *Category
	Handler           CommandHandler
}

type CommandHandler func(*CommandContext)

type CommandContext struct {
	*discordgo.InteractionCreate
	Session *discordgo.Session
}

func (ctx *CommandContext) ReplyWithEmote(emote, message string, args ...interface{}) error {
	return ctx.SendRAW(utils.Fmt("%s | %s", emote, utils.Fmt(message, args...)))
}

func (ctx *CommandContext) SendRAW(message string) error {
	return Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
