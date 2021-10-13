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

func (r *Response) WithRawContent(content string) *Response {
	r.Content = content
	return r
}

func (r *Response) WithContent(emoji, content string, args ...interface{}) *Response {
	return r.WithRawContent(utils.Fmt("%s | %s", emoji, utils.Fmt(content, args...)))
}

func (r *Response) WithFile(file *discordgo.File) *Response {
	r.Files = append(r.Files, file)
	return r
}

func (r *Response) WithRawEmbed(eb *discordgo.MessageEmbed) *Response {
	r.Embeds = append(r.Embeds, eb)
	return r
}

func (r *Response) WithEmbed(embed *Embed) *Response {
	return r.WithRawEmbed(embed.Build())
}

func (r *Response) WithButton(button Button) *Response {
	r.Components = append(r.Components, button.Build())
	return r
}

func (r *Response) ClearComponents() *Response {
	r.Components = []discordgo.MessageComponent{}
	return r
}

func (r *Response) ClearEmbeds() *Response {
	r.Embeds = []*discordgo.MessageEmbed{}
	return r
}

func (r *Response) AsEphemeral() *Response {
	r.Flags = 1 << 6
	return r
}

func (r *Response) Send(channelId string) (*discordgo.Message, error) {
	r.buildBaseComponent()
	return Session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content:         r.Content,
		Files:           r.Files,
		AllowedMentions: r.AllowedMentions,
		Components:      r.Components,
		Embeds:          r.Embeds,
	})
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
