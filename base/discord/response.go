package discord

import (
	"github.com/ItsClairton/Anny/utils"
	"github.com/bwmarrin/discordgo"
)

type Response struct {
	*discordgo.MessageSend
}

func NewResponse() *Response {
	return &Response{&discordgo.MessageSend{
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
	r.Content = utils.Fmt("%s %s", emoji, r.Content)
	return r
}

func (r *Response) WithContentEmoji(emoji, content string, args ...interface{}) *Response {
	r.Content = utils.Fmt("%s %s", emoji, utils.Fmt(content, args...))
	return r
}

func (r *Response) WithFile(file *discordgo.File) *Response {
	r.Files = append(r.Files, file)
	return r
}

func (r *Response) WithEmbed(embed *discordgo.MessageEmbed) *Response {
	r.Embed = embed
	return r
}

func (r *Response) ToInteracctionData() *discordgo.InteractionResponseData {
	data := &discordgo.InteractionResponseData{
		Content:         r.Content,
		Files:           r.Files,
		AllowedMentions: r.AllowedMentions,
	}

	if r.Embed != nil {
		data.Embeds = []*discordgo.MessageEmbed{r.Embed}
	}
	return data
}

func (r *Response) ToWebhookParams() *discordgo.WebhookParams {
	data := &discordgo.WebhookParams{
		Content:         r.Content,
		Files:           r.Files,
		AllowedMentions: r.AllowedMentions,
	}

	if r.Embed != nil {
		data.Embeds = []*discordgo.MessageEmbed{r.Embed}
	}
	return data
}

func (r *Response) ToWebhookEdit() *discordgo.WebhookEdit {
	data := &discordgo.WebhookEdit{
		Content: r.Content,
	}

	if r.Embed != nil {
		data.Embeds = []*discordgo.MessageEmbed{r.Embed}
	}
	return data
}
