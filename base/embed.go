package base

import (
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Embed struct {
	discord.Embed
}

func NewEmbed() *Embed {
	return &Embed{discord.Embed{}}
}

func (e *Embed) SetAuthor(args ...string) *Embed {
	author := &discord.EmbedAuthor{
		Name: args[0],
	}

	if len(args) >= 2 {
		author.Icon = args[1]
	}
	if len(args) >= 3 {
		author.URL = args[2]
	}

	e.Author = author
	return e
}

func (e *Embed) SetTitle(title string, args ...interface{}) *Embed {
	title = utils.Fmt(title, args...)
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

func (e *Embed) SetDescription(desc string, args ...interface{}) *Embed {
	desc = utils.Fmt(desc, args...)
	if len(desc) > 4096 {
		desc = desc[:4096]
	}

	e.Description = desc
	return e
}

func (e *Embed) SetColor(color int) *Embed {
	e.Color = discord.Color(color)
	return e
}

func (e *Embed) SetImage(url string) *Embed {
	e.Image = &discord.EmbedImage{
		URL: url,
	}

	return e
}

func (e *Embed) SetThumbnail(url string) *Embed {
	e.Thumbnail = &discord.EmbedThumbnail{
		URL: url,
	}

	return e
}

func (e *Embed) SetTimestamp(time time.Time) *Embed {
	e.Timestamp = discord.NewTimestamp(time)
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

	e.Fields[index] = discord.EmbedField{
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

	e.Fields = append(e.Fields, discord.EmbedField{
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

	e.Footer = &discord.EmbedFooter{
		Text: content,
		Icon: imgUrl,
	}
	return e
}

func (e *Embed) Build() discord.Embed {
	return e.Embed
}
