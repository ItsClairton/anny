package misc

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PingCommand = base.Interaction{
	Name:        "ping",
	Description: "Pong!",
	Handler: func(ctx *base.InteractionContext) error {
		return ctx.Send(emojis.PingPong, "Pong")
	},
}
