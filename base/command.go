package base

import (
	"strings"

	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/i18n"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name    string
	Aliases []string
	Handler CommandHandler
}

type CommandHandler func(*CommandContext)

type CommandContext struct {
	Message  *discordgo.Message
	Author   *discordgo.User
	Member   *discordgo.Member
	Listener *discordgo.MessageCreate
	Client   *discordgo.Session
	Locale   *i18n.Locale
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

func (ctx *CommandContext) ReplyWithUsage(usage string) (*discordgo.Message, error) {
	return ctx.ReplyWithResponse(response.New(ctx.Locale).WithContentEmote(Emotes.MIKU_CRY, "usage", strings.FieldsFunc(ctx.Message.Content, sutils.SplitString)[0], usage))
}

func (ctx *CommandContext) ReplyWithError(err error) (*discordgo.Message, error) {
	return ctx.ReplyWithResponse(response.New(ctx.Locale).WithContentEmote(Emotes.MIKU_CRY, "error", err.Error()))
}

func (ctx *CommandContext) ReplyWithResponse(response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Listener.ChannelID, response.WithReference(ctx.Message.Reference()).To())
}

func (ctx *CommandContext) EditWithResponse(oldMsg string, response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEditComplex(response.ToEdit(ctx.Message.ChannelID, oldMsg))
}

func (ctx *CommandContext) Reply(emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendReply(ctx.Message.ChannelID, sutils.Fmt("%s | %s", emote, ctx.Locale.GetString(path, args...)), ctx.Message.Reference())
}

func (ctx *CommandContext) ReplyRaw(message string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendReply(ctx.Message.ChannelID, message, ctx.Message.Reference())
}

func (ctx *CommandContext) ReplyRawWithEmote(emote string, message string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendReply(ctx.Message.ChannelID, sutils.Fmt("%s | %s", emote, message), ctx.Message.Reference())
}

func (ctx *CommandContext) SendWithResponse(response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.Message.ChannelID, response.To())
}

func (ctx *CommandContext) SendWithEmote(emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSend(ctx.Message.ChannelID, sutils.Fmt("%s | %s", emote, ctx.Locale.GetString(path, args...)))
}

func (ctx *CommandContext) SendRaw(content string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSend(ctx.Message.ChannelID, content)
}

func (ctx *CommandContext) Send(path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSend(ctx.Message.ChannelID, ctx.Locale.GetString(path, args...))
}

func (ctx *CommandContext) EditWithEmote(id string, emote string, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEdit(ctx.Message.ChannelID, id, sutils.Fmt("%s | %s", emote, ctx.Locale.GetString(path, args...)))
}

func (ctx *CommandContext) DeleteMessage(message *discordgo.Message) {
	ctx.Client.ChannelMessageDelete(message.ChannelID, message.ID)
}
