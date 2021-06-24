package miscellaneous

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/constants"
)

var PingCommand = base.Command{
	Name: "ping",
	Handler: func(ctx *base.CommandContext) {
		ctx.Reply(constants.PING_PONG, "misc.ping.reply", ctx.Client.HeartbeatLatency().Milliseconds())
	},
}
