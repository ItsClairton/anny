package misc

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/sutils"
)

var PingCommand = base.Command{
	Name: "ping", Description: "Saber quantos milissegundos no minimo eu irei responder você",
	Handler: func(ctx *base.CommandContext) {
		ctx.Reply(Emotes.PING_PONG, sutils.Fmt("Pong, `%dms`.", ctx.Client.HeartbeatLatency().Milliseconds()))
	},
}
