package core

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type Command struct {
	Name, Description string
	Module            *Module
	Deffered          bool
	Type              discord.CommandType
	Options           discord.CommandOptions
	Handler           func(*CommandContext)
}

type CommandContext struct {
	*gateway.InteractionCreateEvent

	Data  *discord.CommandInteraction
	State *state.State

	response *api.InteractionResponseData
	sended   bool
}

type Argument struct {
	a *discord.CommandInteractionOption
}

func NewCommandContext(e *gateway.InteractionCreateEvent, state *state.State, data *discord.CommandInteraction, sended bool) *CommandContext {
	return &CommandContext{
		InteractionCreateEvent: e,
		State:                  state,
		Data:                   data,
		response:               &api.InteractionResponseData{},
		sended:                 sended,
	}
}

func (ctx *CommandContext) Argument(index int) *Argument {
	if len(ctx.Data.Options) < index+1 {
		return &Argument{}
	}

	return &Argument{a: &ctx.Data.Options[index]}
}

func (ctx *CommandContext) Guild() *discord.Guild {
	if !ctx.GuildID.IsValid() {
		return nil
	}

	guild, _ := ctx.State.Guild(ctx.GuildID)
	return guild
}

func (ctx *CommandContext) VoiceState() *discord.VoiceState {
	if !ctx.GuildID.IsValid() {
		return nil
	}

	state, _ := ctx.State.VoiceState(ctx.GuildID, ctx.Member.User.ID)
	return state
}

func (ctx *CommandContext) File(file sendpart.File) *CommandContext {
	ctx.response.Files = append(ctx.response.Files, file)
	return ctx
}

func (ctx *CommandContext) Embed(eb *utils.Embed) *CommandContext {
	ctx.response.Embeds = &[]discord.Embed{eb.Embed}
	return ctx
}

func (ctx *CommandContext) Ephemeral() *CommandContext {
	ctx.response.Flags = 1 << 6
	return ctx
}

func (ctx *CommandContext) Reply(args ...interface{}) {
	ctx.checkArguments(args...)

	if ctx.sended {
		ctx.Edit(args...)
		return
	}

	if err := ctx.State.RespondInteraction(ctx.ID, ctx.Token, api.InteractionResponse{Type: 4, Data: ctx.response}); err == nil {
		ctx.sended = true
	} else {
		logger.ErrorF("Não foi possível responder a interação \"%s\" (GuildID: %s): %v", ctx.Data.Name, ctx.GuildID, err)
	}
}

func (ctx *CommandContext) Edit(args ...interface{}) (msg *discord.Message, err error) {
	ctx.checkArguments(args...)

	msg, err = ctx.State.EditInteractionResponse(ctx.AppID, ctx.Token, api.EditInteractionResponseData{
		Content: ctx.response.Content, Components: ctx.response.Components,
		Embeds: ctx.response.Embeds, Files: ctx.response.Files,
	})

	if err != nil {
		logger.ErrorF("Não foi possível editar a resposta da interação \"%s\" (GuildID: %s): %v", ctx.Data.Name, ctx.GuildID, err)
	}

	return
}

func (ctx *CommandContext) checkArguments(args ...interface{}) {
	if len(args) > 1 {
		ctx.response.Content = option.NewNullableString(utils.Fmt("%v | %v", args[0], utils.Fmt(args[1].(string), args[2:]...)))
	}

	if len(args) == 1 {
		if embed, ok := args[0].(*utils.Embed); ok {
			ctx.Embed(embed)
		} else {
			ctx.response.Content = option.NewNullableString(utils.Fmt("%v", args[0]))
		}
	}
}

func (ctx *CommandContext) Stacktrace(err error) {
	if !ctx.sended {
		ctx.Reply(emojis.Cry, "Um erro ocorreu ao executar essa ação: `%v`", err)
	} else {
		ctx.Reply(utils.NewEmbed().Color(0xED4245).Description("%s Um erro ocorreu ao executar essa ação: `%v`", emojis.Cry, err))
	}
}

func (cmd *Command) RAW() api.CreateCommandData {
	return api.CreateCommandData{Name: cmd.Name, Description: cmd.Description, Type: cmd.Type, Options: cmd.Options}
}

func (argument *Argument) Bool() bool {
	if argument.a == nil {
		return false
	}

	value, _ := argument.a.BoolValue()
	return value
}

func (argument *Argument) String() string {
	if argument.a == nil {
		return ""
	}
	return argument.a.String()
}
