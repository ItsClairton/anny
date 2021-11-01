package events

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/bwmarrin/discordgo"
)

var handleFunc = func(i *discord.Interaction, ic *discordgo.InteractionCreate, s *discordgo.Session, sended bool) {
	context := discord.NewContext(ic, s, sended)

	err := i.Handler(context)
	if err != nil {
		logger.Warn(utils.Fmt("Não foi possível responder a interação %s, Guild: %s", i.Name, ic.GuildID), err)
	}

	panic := recover()
	if panic != nil {
		logger.Error(utils.Fmt("Um erro fatal ocorreu ao executar a interação %s, Guild: %s", i.Name, ic.GuildID))
		context.Send(emojis.MikuCry, "Um erro fatal ocorreu ao executar essa ação: `%v`", panic)
	}
}

func InteractionsEvent(s *discordgo.Session, ic *discordgo.InteractionCreate) {

	if ic.Type == discordgo.InteractionApplicationCommand {
		i := discord.GetInteractions()[ic.ApplicationCommandData().Name]
		if i != nil {
			if i.Deffered {
				s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{Type: 5})
				go handleFunc(i, ic, s, true)
			} else {
				go handleFunc(i, ic, s, false)
			}
		}
	}

}
