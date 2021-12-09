package core

import (
	"runtime/debug"

	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func InteractionEvent(e *gateway.InteractionCreateEvent) {
	defer func() {
		if err := recover(); err != nil {
			stacktrace := utils.Fmt("panic: %s\n\n%v", err, string(debug.Stack()))

			State.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString(utils.Fmt("%s | Um erro fatal ocorreu ao executar essa ação: ```go\n%s```", emojis.Cry, stacktrace)),
				},
			})
		}
	}()

	switch data := e.Data.(type) {
	case *discord.CommandInteraction:
		if cmd := Commands[data.Name]; cmd != nil {
			if cmd.Deffered {
				State.RespondInteraction(e.ID, e.Token, api.InteractionResponse{Type: api.DeferredMessageInteractionWithSource})
			}

			cmd.Handler(NewCommandContext(e, State, data, cmd.Deffered))
		}
	}
}
