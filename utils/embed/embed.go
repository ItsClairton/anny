package embed

import "github.com/bwmarrin/discordgo"

type Embed struct {
	*discordgo.MessageEmbed
}

func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

func (e *Embed) Build() *discordgo.MessageEmbed {
	return e.MessageEmbed
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
