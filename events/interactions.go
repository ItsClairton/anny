package events

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var handleFunc = func(i *base.Interaction, ic *gateway.InteractionCreateEvent, s *state.State, sended bool) {
	context := base.NewContext(ic, s, sended)

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

func OnInteraction(e *gateway.InteractionCreateEvent) {
	switch data := e.Data.(type) {
	case *discord.CommandInteraction:
		interaction := base.Interactions[data.Name]
		if interaction != nil {
			if interaction.Deffered {
				base.Session.RespondInteraction(e.ID, e.Token, api.InteractionResponse{Type: 5})
				go handleFunc(interaction, e, base.Session, true)
			} else {
				go handleFunc(interaction, e, base.Session, false)
			}
		}
	}
}
