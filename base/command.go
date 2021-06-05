package base

import (
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name, Description string
	Aliases           []string
	Handler           CommandHandler
}

type CommandHandler func(*CommandContext)

type CommandContext struct {
	Message  *discordgo.Message
	Author   *discordgo.User
	Member   *discordgo.Member
	Listener *discordgo.MessageCreate
	Client   *discordgo.Session
	Args     []string
}

func (ctx *CommandContext) Reply(emote string, message string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendReply(ctx.Message.ChannelID, sutils.Fmt("%s | %s %s", emote, ctx.Author.Mention(), message), ctx.Message.Reference())
}

func (ctx *CommandContext) Send(message string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSend(ctx.Message.ChannelID, message)
}

func (ctx *CommandContext) EditReply(message *discordgo.Message, e string, s string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEdit(message.ChannelID, message.ID, sutils.Fmt("%s | %s %s", e, ctx.Author.Mention(), s))
}

func (ctx *CommandContext) ReplyWithFile(emote string, s string, f *discordgo.File) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Content:   sutils.Fmt("%s | %s %s", emote, ctx.Author.Mention(), s),
		Reference: ctx.Message.Reference(),
		File:      f,
	})
}

func (ctx *CommandContext) DeleteMessage(message *discordgo.Message) {
	ctx.Client.ChannelMessageDelete(message.ChannelID, message.ID)
}

func (ctx *CommandContext) ReplyWithEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Embed:     embed,
		Reference: ctx.Message.Reference(),
	})
}

func (ctx *CommandContext) EditWithEmbed(msg *discordgo.Message, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, embed)
}

func (ctx *CommandContext) ReplyTextWithEmbed(emote string, text string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Message.ID, &discordgo.MessageSend{
		Content:   sutils.Fmt("%s | %s %s", emote, ctx.Author.Mention(), text),
		Embed:     embed,
		Reference: ctx.Message.Reference(),
	})
}
