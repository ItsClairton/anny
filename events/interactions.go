package events

import (
	"runtime/debug"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var handleFunc = func(i *base.Interaction, ic *gateway.InteractionCreateEvent, s *state.State, sended bool) {
	context := base.NewContext(ic, s, sended)

	defer func() {
		if err := recover(); err != nil {
			stacktrace := utils.Fmt("panic: %s\n\n%v", err, string(debug.Stack()))

			logger.ErrorF("Um erro fatal ocorreu ao executar as ações da interação %s, Guilda %s.\n%s", i.Name, ic.GuildID, stacktrace)
			context.SendStackTrace(stacktrace)
		}
	}()

	err := i.Handler(context)
	if err != nil {
		logger.Warn(utils.Fmt("Não foi possível responder a interação %s, Guild: %s", i.Name, ic.GuildID), err)
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
