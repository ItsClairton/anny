package embed

import "github.com/bwmarrin/discordgo"

type Embed struct {
	*discordgo.MessageEmbed
}

func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

func (e *Embed) SetTitle(content string) *Embed {
	if len(content) > 256 {
		content = content[:256]
	}

	e.Title = content
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
