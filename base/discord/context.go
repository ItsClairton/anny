package discord

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

type InteractionContext struct {
	*discordgo.InteractionCreate
	Session  *discordgo.Session
	Response *discordgo.InteractionResponseData
	Sended   bool
}

func NewContext(ic *discordgo.InteractionCreate, s *discordgo.Session, sended bool) *InteractionContext {
	return &InteractionContext{
		InteractionCreate: ic,
		Session:           s,
		Response:          &discordgo.InteractionResponseData{},
		Sended:            sended,
	}
}

func (ctx *InteractionContext) Guild() *discordgo.Guild {
	if ctx.GuildID == "" {
		return nil
	}

	guild, _ := ctx.Session.State.Guild(ctx.GuildID)
	return guild
}

func (ctx *InteractionContext) VoiceState() *discordgo.VoiceState {
	if ctx.GuildID == "" {
		return nil
	}

	state, _ := ctx.Session.State.VoiceState(ctx.GuildID, ctx.Member.User.ID)
	return state
}

func (ctx *InteractionContext) WithContentRAW(content string) *InteractionContext {
	ctx.Response.Content = content
	return ctx
}

func (ctx *InteractionContext) WithContent(emoji, content string, args ...interface{}) *InteractionContext {
	return ctx.WithContentRAW(utils.Fmt("%s | %s", emoji, utils.Fmt(content, args...)))
}

func (ctx *InteractionContext) WithFile(file *discordgo.File) *InteractionContext {
	ctx.Response.Files = append(ctx.Response.Files, file)
	return ctx
}

func (ctx *InteractionContext) WithEmbed(embed *Embed) *InteractionContext {
	ctx.Response.Embeds = append(ctx.Response.Embeds, embed.Build())
	return ctx
}

func (ctx *InteractionContext) AsEphemeral() *InteractionContext {
	ctx.Response.Flags = 1 << 6
	return ctx
}

func (ctx *InteractionContext) Edit(args ...interface{}) error {
	if len(args) > 1 {
		ctx.WithContent(args[0].(string), args[1].(string), args[2:]...)
	}
	if len(args) == 1 {
		ctx.WithContentRAW(args[0].(string))
	}

	_, err := ctx.Session.InteractionResponseEdit(ctx.Session.State.User.ID, ctx.Interaction, &discordgo.WebhookEdit{
		Content: ctx.Response.Content, Components: ctx.Response.Components,
		Embeds: ctx.Response.Embeds, Files: ctx.Response.Files,
	})
	return err
}

func (ctx *InteractionContext) Send(args ...interface{}) error {
	if len(args) > 1 {
		ctx.WithContent(args[0].(string), args[1].(string), args[2:]...)
	}

	if len(args) == 1 {
		ctx.WithContentRAW(args[0].(string))
	}

	if ctx.Sended {
		return ctx.Edit(args...)
	}

	err := ctx.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: ctx.Response,
	})

	if err == nil {
		ctx.Sended = true
	}

	return err
}

func (ctx *InteractionContext) SendWithError(err error) error {
	return ctx.WithEmbed(NewEmbed().
		SetDescription("%s Um erro ocorreu ao executar essa ação: \n```go\n%+v```", emojis.MikuCry, err)).
		Send()
}
