package misc

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var PingCommand = discord.Interaction{
	Name:        "ping",
	Description: "Pong!",
	Handler: func(ctx *discord.InteractionContext) {
		ctx.ReplyWithEmote(emojis.PingPong, "Pong, %dms.", ctx.Session.HeartbeatLatency().Milliseconds())
	},
}
