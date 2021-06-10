package listeners

import (
	"os"
	"strings"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/i18n"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/bwmarrin/discordgo"
)

func MessageCreateListener(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, os.Getenv("DEFAULT_PREFIX")) {
		return
	}

	baseArray := strings.FieldsFunc(m.Content, sutils.SplitString)
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
		Locale:   i18n.GetLocale("pt_BR"),
	})

}
