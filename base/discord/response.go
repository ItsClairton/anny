package discord

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/bwmarrin/discordgo"
)

type Response struct {
	*discordgo.InteractionResponseData
}

func NewResponse() *Response {
	return &Response{&discordgo.InteractionResponseData{
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	}}
}

func (r *Response) WithContent(content string, args ...interface{}) *Response {
	r.Content = utils.Fmt(content, args...)
	return r
}

func (r *Response) WithEmoji(emoji string) *Response {
	r.Content = utils.Fmt("%s | %s", emoji, r.Content)
	return r
}

func (r *Response) WithContentEmoji(emoji, content string, args ...interface{}) *Response {
	r.Content = utils.Fmt("%s | %s", emoji, utils.Fmt(content, args...))
	return r
}

func (r *Response) WithFile(file *discordgo.File) *Response {
	r.Files = append(r.Files, file)
	return r
}

func (r *Response) WithEmbed(embed *discordgo.MessageEmbed) *Response {
	r.Embeds = append(r.Embeds, embed)
	return r
}

func (r *Response) WithButton(button Button) *Response {
	r.Components = append(r.Components, button.Build())
	return r
}

func (r *Response) ClearComponents() *Response {
	r.Components = []discordgo.MessageComponent{}
	return r
}

func (r *Response) AsEphemeral() *Response {
	r.Flags = 1 << 6
	return r
}

func (r *Response) Build() *discordgo.InteractionResponseData {
	return r.buildBaseComponent().InteractionResponseData
}

func (r *Response) BuildAsWebhookParams() *discordgo.WebhookParams {
	r.buildBaseComponent()

	return &discordgo.WebhookParams{
		Content:         r.Content,
		Files:           r.Files,
		AllowedMentions: r.AllowedMentions,
		Components:      r.Components,
		Embeds:          r.Embeds,
		Flags:           r.Flags,
	}
}

func (r *Response) BuildAsWebhookEdit() *discordgo.WebhookEdit {
	r.buildBaseComponent()

	return &discordgo.WebhookEdit{
		Content:    r.Content,
		Components: r.Components,
		Embeds:     r.Embeds,
		Files:      r.Files,
	}
}

// TODO: Ver uma forma melhor de fazer isso
func (r *Response) buildBaseComponent() *Response {
	if len(r.Components) == 0 {
		return r
	}

	r.Components = []discordgo.MessageComponent{discordgo.ActionsRow{
		Components: r.Components,
	}}
	return r
}
