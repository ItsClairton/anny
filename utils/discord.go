package utils

import "github.com/bwmarrin/discordgo"

func GetFirstAttachment(msg *discordgo.Message) string {

	if len(msg.Attachments) > 0 {
		return msg.Attachments[0].ProxyURL
	}

	if len(msg.Embeds) > 0 && (msg.Embeds[0].Image != nil || msg.Embeds[0].Thumbnail != nil) {
		if msg.Embeds[0].Image != nil {
			return msg.Embeds[0].Image.ProxyURL
		}
		return msg.Embeds[0].Thumbnail.ProxyURL
	}

	return ""
}
