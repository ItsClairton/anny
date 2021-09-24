package discord

import "github.com/bwmarrin/discordgo"

type Embed struct {
	*discordgo.MessageEmbed
}

func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
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

func (e *Embed) SetTitle(title string) *Embed {
	if len(title) > 256 {
		title = title[:256]
	}

	e.Title = title
	return e
}

func (e *Embed) SetURL(url string) *Embed {
	e.URL = url
	return e
}

func (e *Embed) SetDescription(desc string) *Embed {
	if len(desc) > 4096 {
		desc = desc[:4096]
	}

	e.Description = desc
	return e
}

func (e *Embed) SetColor(color int) *Embed {
	e.Color = color
	return e
}

func (e *Embed) SetImage(url string) *Embed {
	e.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}

	return e
}

func (e *Embed) SetThumbnail(url string) *Embed {
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}

	return e
}

func (e *Embed) SetField(index int, name, value string, inline bool) *Embed {
	if index > len(e.Fields) {
		return e
	}
	if len(name) > 256 {
		name = name[:256]
	}
	if len(value) > 1024 {
		value = value[:1024]
	}

	e.Fields[index] = &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	return e
}

func (e *Embed) AddField(name, value string, inline bool) *Embed {
	if len(name) > 256 {
		name = name[:256]
	}
	if len(value) > 1024 {
		value = value[:1024]
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return e
}

func (e *Embed) SetFooter(content, imgUrl string) *Embed {
	if len(content) > 2048 {
		content = content[:2048]
	}

	e.Footer = &discordgo.MessageEmbedFooter{
		Text:    content,
		IconURL: imgUrl,
	}
	return e
}

func (e *Embed) Build() *discordgo.MessageEmbed {
	return e.MessageEmbed
}
