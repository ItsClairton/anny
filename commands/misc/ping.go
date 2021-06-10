package misc

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/Emotes"
)

var PingCommand = base.Command{
	Name: "ping",
	Handler: func(ctx *base.CommandContext) {
		ctx.Reply(Emotes.PING_PONG, "misc.ping.reply", ctx.Client.HeartbeatLatency().Milliseconds())
	},
}
