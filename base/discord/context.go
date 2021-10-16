package discord

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

type InteractionContext struct {
	Session       *discordgo.Session
	ResponseType  int
	AlreadySended bool
	*discordgo.InteractionCreate
}

func (ctx *InteractionContext) GetGuild() (*discordgo.Guild, error) {
	if ctx.GuildID == "" {
		return nil, nil
	}

	return ctx.Session.State.Guild(ctx.GuildID)
}

func (ctx *InteractionContext) GetVoiceChannel() string {
	guild, err := ctx.GetGuild()
	if guild == nil || err != nil {
		return ""
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Member.User.ID {
			return vs.ChannelID
		}
	}
	return ""
}

func (ctx *InteractionContext) SendComplex(data *discordgo.InteractionResponseData) (*discordgo.Message, error) {
	err := ctx.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(ctx.ResponseType),
		Data: data,
	})
	if err == nil {
		ctx.AlreadySended = true
	}

	return nil, err
}

func (ctx *InteractionContext) EditComplex(data *discordgo.WebhookEdit) (*discordgo.Message, error) {
	return ctx.Session.InteractionResponseEdit(ctx.Session.State.User.ID, ctx.Interaction, data)
}

func (ctx *InteractionContext) SendResponse(res *Response) (*discordgo.Message, error) {
	if ctx.AlreadySended {
		return ctx.EditComplex(res.BuildAsWebhookEdit())
	}

	return ctx.SendComplex(res.Build())
}

func (ctx *InteractionContext) SendEmbed(embeds ...*discordgo.MessageEmbed) (*discordgo.Message, error) {
	if ctx.AlreadySended {
		return ctx.EditComplex(&discordgo.WebhookEdit{
			Embeds:     embeds,
			Components: []discordgo.MessageComponent{},
		})
	}

	return ctx.SendComplex(&discordgo.InteractionResponseData{
		Embeds: embeds,
	})
}

func (ctx *InteractionContext) SendFile(files ...*discordgo.File) (*discordgo.Message, error) {
	if ctx.AlreadySended {
		return ctx.EditComplex(&discordgo.WebhookEdit{
			Files: files,
		})
	}

	return ctx.SendComplex(&discordgo.InteractionResponseData{
		Files: files,
	})
}

func (ctx *InteractionContext) SendRAW(text string) (*discordgo.Message, error) {
	if ctx.AlreadySended {
		return ctx.EditComplex(&discordgo.WebhookEdit{
			Content: text,
		})
	}

	return ctx.SendComplex(&discordgo.InteractionResponseData{
		Content: text,
	})
}

func (ctx *InteractionContext) SendError(err error) {
	ctx.SendEmbed(NewEmbed().
		SetColor(0xF93A2F).
		SetDescription(utils.Fmt("%s Um erro ocorreu ao executar essa ação: `%s`", emojis.MikuCry, err.Error())).
		Build())
}

func (ctx *InteractionContext) SendEphemeral(emoji, text string, args ...interface{}) (*discordgo.Message, error) {
	return ctx.SendComplex(&discordgo.InteractionResponseData{
		Content: utils.Fmt("%s | %s", emoji, utils.Fmt(text, args...)),
		Flags:   1 << 6,
	})
}

func (ctx *InteractionContext) Send(emoji, text string, args ...interface{}) (*discordgo.Message, error) {
	content := utils.Fmt("%s | %s", emoji, utils.Fmt(text, args...))
	if ctx.AlreadySended {
		return ctx.EditComplex(&discordgo.WebhookEdit{
			Content:    content,
			Components: []discordgo.MessageComponent{},
		})
	}

	return ctx.SendRAW(content)
}
