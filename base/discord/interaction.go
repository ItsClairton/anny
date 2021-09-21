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
	Session *discordgo.Session
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

func (ctx *InteractionContext) ReplyWithEmote(emote, message string, args ...interface{}) error {
	return ctx.SendRAW(utils.Fmt("%s | %s", emote, utils.Fmt(message, args...)))
}

func (ctx *InteractionContext) ReplyEphemeralWithEmote(emote, message string, args ...interface{}) error {
	return ctx.SendEphemeralRAW(utils.Fmt("%s | %s", emote, utils.Fmt(message, args...)))
}

func (ctx *InteractionContext) SendError(err error) {
	ctx.ReplyEphemeralWithEmote(emojis.MikuCry, "Um erro ocorreu ao executar esse comando, Desculpa (`%v`).", err)
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
