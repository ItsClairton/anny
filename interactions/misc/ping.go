package misc

import (
	"time"

	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PingCommand = core.Interaction{
	Name:        "ping",
	Description: "Pong!",
	Handler: func(ctx *core.InteractionContext) error {
		latency := time.Duration(ctx.Session.PacerLoop.EchoBeat.Get() - ctx.Session.PacerLoop.SentBeat.Get())
		if latency == 0 {
			return ctx.Send(emojis.PingPong, "Não há métricas de latência ainda ;(")
		}

		return ctx.Send(emojis.PingPong, "Pong, %dms.", latency.Milliseconds())
	},
}
