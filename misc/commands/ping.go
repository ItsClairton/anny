package misc

import (
	"time"

	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PingCommand = core.Command{
	Name:        "ping",
	Description: "Mostra uma média da latência do bot para os servidores do Discord",
	Handler: func(ctx *core.CommandContext) {
		latency := time.Duration(ctx.State.PacerLoop.EchoBeat.Get() - ctx.State.PacerLoop.SentBeat.Get())
		if latency <= 0 {
			ctx.Reply(emojis.PingPong, "Não há medições de latência suficientes ainda ;(")
			return
		}

		ctx.Reply(emojis.PingPong, "Pong, %dms.", latency.Milliseconds())
	},
}
