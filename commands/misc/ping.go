package misc

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PingCommand = discord.Command{
	Name:        "ping",
	Description: "Pong!",
	Handler: func(ctx *discord.CommandContext) {
		ctx.ReplyWithEmote(emojis.PingPong, "Pong, %dms.", ctx.Session.HeartbeatLatency().Milliseconds())
	},
}
