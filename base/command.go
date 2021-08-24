package base

import (
	"strings"

	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/constants"
	"github.com/ItsClairton/Anny/utils/i18n"
	"github.com/bwmarrin/discordgo"
)

type Category struct {
	ID, Emote string
	Commands  []*Command
}

type Command struct {
	Name     string
	Category *Category
	Aliases  []string
	Handler  CommandHandler
}

type CommandHandler func(*CommandContext)

type CommandContext struct {
	*i18n.Locale
	*discordgo.MessageCreate
	Client *discordgo.Session
	Args   []string
}

func (ctx *CommandContext) GetGuild() *discordgo.Guild {
	g, _ := ctx.Client.State.Guild(ctx.GuildID)
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

func (ctx *CommandContext) GetArgsWithLines() []string {

	args := strings.Split(ctx.Message.Content, " ")

	if len(args) > 1 {
		return args[1:]
	} else {
		return nil
	}

}

func (ctx *CommandContext) Reply(emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.ReplyRaw(utils.Fmt("%s | %s", emote, ctx.GetString(path, args...)))
}

func (ctx *CommandContext) ReplyWithoutEmote(path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.ReplyRaw(ctx.GetString(path, args...))
}

func (ctx *CommandContext) ReplyRaw(message string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.ChannelID, &discordgo.MessageSend{
		Content:   message,
		Reference: ctx.Reference(),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
}

func (ctx *CommandContext) ReplyRawWithEmote(emote, message string) (*discordgo.Message, error) {
	return ctx.ReplyRaw(utils.Fmt("%s | %s", emote, message))
}

func (ctx *CommandContext) ReplyWithResponse(response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.ChannelID, response.WithReference(ctx.Reference()).To())
}

func (ctx *CommandContext) ReplyWithEmbed(eb *embed.Embed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.ChannelID, &discordgo.MessageSend{
		Embed:     eb.Build(),
		Reference: ctx.Reference(),
	})
}

func (ctx *CommandContext) ReplyWithUsage(usage string) (*discordgo.Message, error) {
	return ctx.Reply(constants.MIKU_CRY, "usage", strings.FieldsFunc(ctx.Message.Content, utils.SplitString)[0], usage)
}

func (ctx *CommandContext) ReplyWithError(err error) (*discordgo.Message, error) {
	return ctx.Reply(constants.MIKU_CRY, "error", err.Error())
}

func (ctx *CommandContext) Send(emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.SendRaw(utils.Fmt("%s | %s", ctx.GetString(path, args...)))
}

func (ctx *CommandContext) SendWithoutEmote(path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.SendRaw(ctx.GetString(path, args...))
}

func (ctx *CommandContext) SendRaw(content string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.ChannelID, &discordgo.MessageSend{
		Content: content,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
}

func (ctx *CommandContext) SendWithEmbed(eb *embed.Embed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendEmbed(ctx.ChannelID, eb.Build())
}

func (ctx *CommandContext) SendWithEmbedTo(channelId string, eb *embed.Embed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendEmbed(channelId, eb.Build())
}

func (ctx *CommandContext) SendWithResponse(response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageSendComplex(ctx.ChannelID, response.To())
}

func (ctx *CommandContext) Edit(msgId, emote, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.EditRaw(msgId, utils.Fmt("%s | %s", ctx.GetString(path, args...)))
}

func (ctx *CommandContext) EditWithoutEmote(msgId, path string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.EditRaw(msgId, ctx.GetString(path, args...))
}

func (ctx *CommandContext) EditRaw(msgId string, content string) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      msgId,
		Channel: ctx.ChannelID,
		Content: &content,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
}

func (ctx *CommandContext) EditWithEmbed(msgId string, eb *embed.Embed) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEditEmbed(ctx.ChannelID, msgId, eb.Build())
}

func (ctx *CommandContext) EditWithResponse(msgId string, response *response.Response) (*discordgo.Message, error) {
	return ctx.Client.ChannelMessageEditComplex(response.ToEdit(ctx.ChannelID, msgId))
}

func (ctx *CommandContext) DeleteMessage(message *discordgo.Message) {
	ctx.Client.ChannelMessageDelete(message.ChannelID, message.ID)
}
