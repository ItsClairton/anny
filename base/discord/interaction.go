package discord

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

type Category struct {
	Name, Emote  string
	Interactions []*Interaction
}

type Interaction struct {
	Name, Description string
	Type              discordgo.ApplicationCommandType
	Options           []*discordgo.ApplicationCommandOption
	Category          *Category
	Handler           InteractionHandler
}

type InteractionContext struct {
	*discordgo.InteractionCreate
	Session  *discordgo.Session
	Deffered bool
}

type InteractionHandler func(*InteractionContext)

func (i Interaction) ToRAW() *discordgo.ApplicationCommand {
	raw := &discordgo.ApplicationCommand{
		Name:        i.Name,
		Description: i.Description,
		Type:        i.Type,
		Options:     i.Options,
	}

	return raw
}

func (ctx *InteractionContext) SendDeffered(ephemeral bool) error {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}

	if ephemeral {
		response.Data = &discordgo.InteractionResponseData{
			Flags: 1 << 6,
		}
	}
	err := ctx.Session.InteractionRespond(ctx.Interaction, response)
	if err == nil {
		ctx.Deffered = true
	}

	return err
}

func (ctx *InteractionContext) ReplyWithEmote(emote, message string, args ...interface{}) error {
	return ctx.SendRAW(utils.Fmt("%s | %s", emote, utils.Fmt(message, args...)))
}

func (ctx *InteractionContext) ReplyEphemeralWithEmote(emote, message string, args ...interface{}) error {
	return ctx.SendEphemeralRAW(utils.Fmt("%s | %s", emote, utils.Fmt(message, args...)))
}

func (ctx *InteractionContext) SendError(err error) {
	if ctx.Deffered {
		ctx.EditWithEmote(emojis.MikuCry, "Um erro ocorreu ao executar esse comando: `%s`", err.Error())
	} else {
		ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Um erro ocorreu ao executar esse comando: `%s`", err.Error())
	}
}

func (ctx *InteractionContext) SendEphemeralRAW(message string) error {
	return Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   1 << 6,
		},
	})
}

func (ctx *InteractionContext) SendRAW(message string) error {
	return Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func (ctx *InteractionContext) SendResponse(response *Response) error {
	return Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: response.Build(),
	})
}

func (ctx *InteractionContext) EditWithEmote(emote, message string, args ...interface{}) (*discordgo.Message, error) {
	return Session.InteractionResponseEdit(Session.State.User.ID, ctx.Interaction, &discordgo.WebhookEdit{
		Content: utils.Fmt("%s | %s", emote, utils.Fmt(message, args...)),
	})
}

func (ctx *InteractionContext) EditResponse(response *Response) (*discordgo.Message, error) {
	return Session.InteractionResponseEdit(ctx.Session.State.User.ID, ctx.Interaction, response.BuildAsWebhookEdit())
}

func (ctx *InteractionContext) DeleteResponse() error {
	return Session.InteractionResponseDelete(ctx.Session.State.User.ID, ctx.Interaction)
}

func (ctx *InteractionContext) SendFollowUp(response *Response) (*discordgo.Message, error) {
	return Session.FollowupMessageCreate(ctx.Session.State.User.ID, ctx.Interaction, true, response.BuildAsWebhookParams())
}

func (ctx *InteractionContext) EditFollowUp(id string, response *Response) (*discordgo.Message, error) {
	return Session.FollowupMessageEdit(ctx.Session.State.User.ID, ctx.Interaction, id, response.BuildAsWebhookEdit())
}

func (ctx *InteractionContext) DeleteFollowUp(id string) error {
	return Session.FollowupMessageDelete(ctx.Session.State.User.ID, ctx.Interaction, id)
}

func (ctx *InteractionContext) GetGuild() *discordgo.Guild {
	guild, err := ctx.Session.State.Guild(ctx.GuildID)

	if err != nil {
		return nil
	}
	return guild
}

func (ctx *InteractionContext) GetVoiceChannel() string {
	guild := ctx.GetGuild()
	if guild == nil {
		return ""
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Member.User.ID {
			return vs.ChannelID
		}
	}
	return ""
}
