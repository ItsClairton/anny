package base

import (
	"strings"

	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/i18n"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/constants"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name    string
	Aliases []string
	Handler CommandHandler
}

type CommandHandler func(*CommandContext)

type CommandContext struct {
	*i18n.Locale
	Message  *discordgo.Message
	Author   *discordgo.User
	Member   *discordgo.Member
	Listener *discordgo.MessageCreate
	Client   *discordgo.Session
	Args     []string
}

func (ctx *CommandContext) GetGuild() *discordgo.Guild {
	g, _ := ctx.Client.State.Guild(ctx.Message.GuildID)
	return g
}

func (ctx *CommandContext) GetVoice() string {
	for _, vs := range ctx.GetGuild().VoiceStates {
		if vs.UserID == ctx.Author.ID {
			return vs.ChannelID
		}
	}

	return ""
}

func (ctx *CommandContext) Reply(emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.ReplyRaw(utils.Fmt("%s | %s", emote, ctx.Locale.GetString(path, args...)))
}

func (ctx *CommandContext) ReplyWithoutEmote(path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.ReplyRaw(ctx.Locale.GetString(path, args...))
}

func (ctx *CommandContext) ReplyRaw(message string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendReply(ctx.Message.ChannelID, message, ctx.Message.Reference())
}

func (ctx *CommandContext) ReplyRawWithEmote(emote, message string) (*discordgo.Message, error) {
	return ctx.ReplyRaw(utils.Fmt("%s | %s", emote, message))
}

func (ctx *CommandContext) ReplyWithResponse(response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Listener.ChannelID, response.WithReference(ctx.Message.Reference()).To())
}

func (ctx *CommandContext) ReplyWithEmbed(eb *embed.Embed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Listener.ChannelID, &discordgo.MessageSend{
		Embed:     eb.Build(),
		Reference: ctx.Message.Reference(),
	})
}

func (ctx *CommandContext) ReplyWithUsage(usage string) (*discordgo.Message, error) {
	return ctx.Reply(constants.MIKU_CRY, "usage", strings.FieldsFunc(ctx.Message.Content, utils.SplitString)[0], usage)
}

func (ctx *CommandContext) ReplyWithError(err error) (*discordgo.Message, error) {
	return ctx.Reply(constants.MIKU_CRY, "error", err.Error())
}

func (ctx *CommandContext) Send(emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.SendRaw(utils.Fmt("%s | %s", ctx.Locale.GetString(path, args...)))
}

func (ctx *CommandContext) SendWithoutEmote(path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.SendRaw(ctx.Locale.GetString(path, args...))
}

func (ctx *CommandContext) SendRaw(content string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSend(ctx.Message.ChannelID, content)
}

func (ctx *CommandContext) SendWithEmbed(eb *embed.Embed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendEmbed(ctx.Message.ChannelID, eb.Build())
}

func (ctx *CommandContext) SendWithResponse(response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Message.ChannelID, response.To())
}

func (ctx *CommandContext) Edit(msgId, emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.EditRaw(msgId, utils.Fmt("%s | %s", ctx.Locale.GetString(path, args...)))
}

func (ctx *CommandContext) EditWithoutEmote(msgId, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.EditRaw(msgId, ctx.Locale.GetString(path, args...))
}

func (ctx *CommandContext) EditRaw(msgId string, content string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEdit(ctx.Listener.ChannelID, msgId, content)
}

func (ctx *CommandContext) EditWithEmbed(msgId string, eb *embed.Embed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEditEmbed(ctx.Listener.ChannelID, msgId, eb.Build())
}

func (ctx *CommandContext) EditWithResponse(msgId string, response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEditComplex(response.ToEdit(ctx.Message.ChannelID, msgId))
}

func (ctx *CommandContext) DeleteMessage(message *discordgo.Message) {
	ctx.Client.ChannelMessageDelete(message.ChannelID, message.ID)
}
