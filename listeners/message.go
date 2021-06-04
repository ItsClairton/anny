package listeners

import (
	"os"
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/bwmarrin/discordgo"
)

func splitString(r rune) bool {
	return r == ' ' || r == '\n'
}

func MessageCreateListener(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, os.Getenv("DEFAULT_PREFIX")) {
		return
	}

	baseArray := strings.FieldsFunc(m.Content, splitString)
	label := strings.ToLower(strings.TrimPrefix(baseArray[0], os.Getenv("DEFAULT_PREFIX")))

	cmd, exist := base.GetCommandMapper()[label]

	if !exist {
		return
	}

	go s.ChannelTyping(m.ChannelID)

	var args []string

	if len(baseArray) > 1 {
		args = baseArray[1:]
	}

	go cmd.Handler(&base.CommandContext{
		Message:  m.Message,
		Author:   m.Author,
		Listener: m,
		Member:   m.Member,
		Client:   s,
		Args:     args,
	})

}
