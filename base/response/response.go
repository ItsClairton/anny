package response

import (
	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/i18n"
	"github.com/ItsClairton/Anny/utils"
	"github.com/bwmarrin/discordgo"
)

type Response struct {
	locale *i18n.Locale
	*discordgo.MessageSend
}

func New(locale *i18n.Locale) *Response {
	return &Response{locale, &discordgo.MessageSend{
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	}}
}

func (r *Response) ClearContent() *Response {
	r.Content = ""
	return r
}

func (r *Response) SetContentEmote(emote, content string) *Response {
	r.Content = utils.Fmt("%s | %s", emote, content)
	return r
}

func (r *Response) WithContentEmote(emote, path string, args ...interface{}) *Response {
	r.Content = utils.Fmt("%s | %s", emote, r.locale.GetString(path, args...))
	return r
}

func (r *Response) WithReference(ref *discordgo.MessageReference) *Response {
	r.Reference = ref
	return r
}

func (r *Response) WithContent(path string, content ...interface{}) *Response {
	r.Content = r.locale.GetString(path, content...)
	return r
}

func (r *Response) WithFile(file *discordgo.File) *Response {
	r.Files = append(r.Files, file)
	return r
}

func (r *Response) WithEmbed(eb *embed.Embed) *Response {
	r.Embed = eb.Build()
	return r
}

func (r *Response) To() *discordgo.MessageSend {
	return r.MessageSend
}

func (r *Response) ToEdit(ch string, id string) *discordgo.MessageEdit {
	return &discordgo.MessageEdit{
		Channel:         ch,
		ID:              id,
		Content:         &r.Content,
		Embed:           r.Embed,
		AllowedMentions: r.AllowedMentions,
	}
}
