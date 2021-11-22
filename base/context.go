package base

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type InteractionContext struct {
	*gateway.InteractionCreateEvent
	Session  *state.State
	Response *api.InteractionResponseData
	Sended   bool
}

func NewContext(ic *gateway.InteractionCreateEvent, s *state.State, sended bool) *InteractionContext {
	return &InteractionContext{
		InteractionCreateEvent: ic,
		Session:                s,
		Response:               &api.InteractionResponseData{},
		Sended:                 sended,
	}
}

func (ctx *InteractionContext) DataAsCommand() *discord.CommandInteraction {
	if ctx.Data.InteractionType() != discord.CommandInteractionType {
		return nil
	}

	return ctx.Data.(*discord.CommandInteraction)
}

func (ctx *InteractionContext) Argument(index int) (discord.CommandInteractionOption, bool) {
	if data := ctx.DataAsCommand(); data == nil || len(data.Options) < index+1 {
		return discord.CommandInteractionOption{}, false
	} else {
		return data.Options[index], true
	}
}

func (ctx *InteractionContext) ArgumentAsBool(index int) bool {
	argument, exist := ctx.Argument(index)
	if !exist {
		return false
	}

	parsed, _ := argument.BoolValue()
	return parsed
}

func (ctx *InteractionContext) ArgumentAsString(index int) string {
	argument, _ := ctx.Argument(index)
	return argument.String()
}

func (ctx *InteractionContext) ArgumentAsInteger(index int) int64 {
	argument, exist := ctx.Argument(index)
	if !exist {
		return -1
	}

	parsed, _ := argument.IntValue()
	return parsed
}

func (ctx *InteractionContext) Guild() *discord.Guild {
	if !ctx.GuildID.IsValid() {
		return nil
	}

	guild, _ := ctx.Session.Guild(ctx.GuildID)
	return guild
}

func (ctx *InteractionContext) VoiceState() *discord.VoiceState {
	if !ctx.GuildID.IsValid() {
		return nil
	}

	state, _ := ctx.Session.VoiceState(ctx.GuildID, ctx.Member.User.ID)

	return state
}

func (ctx *InteractionContext) WithContentRAW(content string) *InteractionContext {
	ctx.Response.Content = option.NewNullableString(content)
	return ctx
}

func (ctx *InteractionContext) WithContent(emoji, content string, args ...interface{}) *InteractionContext {
	return ctx.WithContentRAW(utils.Fmt("%s | %s", emoji, utils.Fmt(content, args...)))
}

func (ctx *InteractionContext) WithFile(file sendpart.File) *InteractionContext {
	ctx.Response.Files = append(ctx.Response.Files, file)
	return ctx
}

func (ctx *InteractionContext) WithEmbed(embed *Embed) *InteractionContext {
	ctx.Response.Embeds = &[]discord.Embed{embed.Build()}

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

	_, err := ctx.Session.EditInteractionResponse(ctx.AppID, ctx.Token, api.EditInteractionResponseData{
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

	err := ctx.Session.RespondInteraction(ctx.ID, ctx.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: ctx.Response,
	})

	if err == nil {
		ctx.Sended = true
	}

	return err
}

func (ctx *InteractionContext) SendError(err error) error {
	return ctx.WithEmbed(NewEmbed().SetColor(0xE74C3C).SetDescription("%s Um erro ocorreu ao executar essa ação: `%v`", emojis.Cry, err)).Send()
}

func (ctx *InteractionContext) SendStackTrace(stack string) error {
	return ctx.WithEmbed(NewEmbed().SetColor(0xE74C3C).
		SetTitle("%s Um erro ocorreu ao executar essa ação:", emojis.Cry).
		SetDescription("```go\n%+v```", stack)).
		Send()
}
