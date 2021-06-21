package embed

import (
	"github.com/ItsClairton/Anny/utils/i18n"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/bwmarrin/discordgo"
)

type Embed struct {
	*discordgo.MessageEmbed
	key    string
	locale *i18n.Locale
}

func NewEmbed(locale *i18n.Locale, path string) *Embed {
	return &Embed{&discordgo.MessageEmbed{}, path, locale}
}

func (e *Embed) Build() *discordgo.MessageEmbed {
	return e.MessageEmbed
}

func (e *Embed) WithAuthor(iconUrl, url string, args ...interface{}) *Embed {
	return e.SetAuthor(e.locale.GetString(sutils.Fmt("%s.author", e.key), args...), iconUrl, url)
}

func (e *Embed) SetAuthor(args ...string) *Embed {

	author := &discordgo.MessageEmbedAuthor{
		Name: args[0],
	}

	if len(args) >= 2 {
		author.IconURL = args[1]
	}

	if len(args) >= 3 {
		author.URL = args[2]
	}

	e.Author = author
	return e
}

func (e *Embed) WithTitle(args ...interface{}) *Embed {
	e.SetTitle(e.locale.GetString(sutils.Fmt("%s.title", e.key), args...))
	return e
}

func (e *Embed) SetTitle(content string) *Embed {
	if len(content) > 256 {
		content = content[:256]
	}

	e.Title = content
	return e
}

func (e *Embed) WithDescription(args ...interface{}) *Embed {
	e.SetDescription(e.locale.GetString(sutils.Fmt("%s.description", e.key), args...))
	return e
}

func (e *Embed) WithEmoteDescription(emote string, args ...interface{}) *Embed {
	e.SetDescription(sutils.Fmt("%s %s", emote, e.locale.GetString(sutils.Fmt("%s.description", e.key), args...)))
	return e
}

func (e *Embed) SetDescription(content string) *Embed {
	if len(content) > 2048 {
		content = content[:2048]
	}

	e.Description = content
	return e
}

func (e *Embed) SetColor(color int) *Embed {
	e.Color = color
	return e
}

func (e *Embed) SetURL(url string) *Embed {
	e.URL = url
	return e
}

func (e *Embed) SetThumbnail(url string) *Embed {
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	return e
}

func (e *Embed) SetImage(url string) *Embed {
	e.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return e
}

func (e *Embed) WithField(value string, inline bool) *Embed {
	return e.AddField(e.locale.GetString(sutils.Fmt("%s.fields.[%v]", e.key, len(e.Fields))), value, inline)
}

func (e *Embed) SetFieldValue(index int, value string) *Embed {
	if index > len(e.Fields) {
		return e
	}

	if len(value) > 1024 {
		value = value[:1024]
	}
	if len(value) < 1 {
		value = "N/A"
	}

	e.Fields[index].Value = value
	return e
}

func (e *Embed) SetField(index int, title string, value string, inline bool) *Embed {
	if index > len(e.Fields) {
		return e
	}

	if len(title) > 256 {
		title = title[:256]
	}

	if len(value) > 1024 {
		value = value[:1024]
	}

	if len(value) < 1 {
		value = "N/A"
	}

	e.Fields[index] = &discordgo.MessageEmbedField{
		Name:   title,
		Value:  value,
		Inline: inline,
	}

	return e
}

func (e *Embed) AddField(title, value string, inline bool) *Embed {

	if len(title) > 256 {
		title = title[:256]
	}

	if len(value) > 1024 {
		value = value[:1024]
	}

	if len(value) < 1 {
		value = "N/A"
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   title,
		Value:  value,
		Inline: inline,
	})
	return e
}

func (e *Embed) SetFooter(content string, imgUrl string) *Embed {

	e.Footer = &discordgo.MessageEmbedFooter{
		Text:    content,
		IconURL: imgUrl,
	}
	return e
}

func (e *Embed) WithFooter(imgUrl string, values ...interface{}) *Embed {
	return e.SetFooter(imgUrl, e.locale.GetString(sutils.Fmt("%s.footer", e.key), values...))
}
